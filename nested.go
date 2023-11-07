// Вложенные объекты в Go.
//
// Используются как альтернатива словарям с динамической типизацией для реализации работы с JSON-объектами.
//
// Не являются конкурентно-безопасными, для добавления этой возможности надо делать обертку.
//
// Пример использования для создания и чтения объекта:
//
//	// {
//	// 	 "nested_object": {
//	// 	   "key1": "value1",
//	//	   "key2": "value2"
//	//	 },
//	//	 "nested_array": ["elem1", "elem2"]
//	// }
//
//	nested := Nested{}
//
//	nested.SetValue("value1", "nested_object", "key1")
//	nested.SetValue("value2", "nested_object", "key2")
//	nested.SetArray([]*Nested{}, "nested_array")
//	nested.ArrayAddValue("elem1", "nested_array")
//	nested.ArrayAddValue("elem2", "nested_array")
//
//	nested.GetValue("nested_object", "key1")
//	nested.GetValue("nested_object", "key2")
//	array, _ := nested.GetArray("nested_array")
//	for _, elem := range array {
//		elem.GetValue()
//	}
//
// Eсть возможность инициализации структуры из JSON-строки и обратно с помощью [FromJSONString] и [ToJSONString]:
//
//	nested := FromJSONString(`{"nested_object": {"key1": "value1", "key2": "value2"}, "nested_array": ["elem1", "elem2"]}`)
//
//	nested.GetValue("nested_object", "key1")
//	nested.GetValue("nested_object", "key2")
//	array, _ := nested.GetArray("nested_array")
//	for _, elem := range array {
//		elem.GetValue()
//	}
//
//	nested.ToJSONString() // {"nested_object": {"key1": "value1", "key2": "value2"}, "nested_array": ["elem1", "elem2"]}
//
// Для инициализации структуры словарем, массивом или значением-интерфейсом
// и обратной конвертации в интерфейс см. [FromObject] и [ToObject].
package nested

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// Структура для описания объекта.
//
// Объект может быть скалярным значением любого типа (any), массивом указателей на объекты или объектом вида ключ-значение,
// где ключ - строка, значение - указатель на объект.
//
// Для каждого объекта может использоваться только один из этих вариантов поведения.
// По умолчанию объект имеет вид ключ-значение. Другие варианты обозначаются выставлением флагов isValue и isArray.
//
// Поля в объекте не экспортируются, чтобы избежать случайных конфликтов при инициализации.
type Nested struct {
	isValue bool // является ли объект скалярным значением
	isArray bool // является ли объект массивом

	nested map[string]*Nested // вложенный объект вида ключ-значение
	array  []*Nested          // массив объектов

	value any // скалярное значение
}

// Проверка, является ли объект массивом.
func (j *Nested) IsArray() bool {
	return j.isArray
}

// Проверка, является ли объект скалярным значением.
func (j *Nested) IsValue() bool {
	return j.isValue
}

// Проверка, что объект имеет тип ключ-значение.
func (j *Nested) IsNested() bool {
	return !j.isArray && !j.isValue
}

// Проверка, что объект пустой и имеет вид ключ-значение.
func (j *Nested) IsEmpty() bool {
	return !j.isArray && !j.isValue && len(j.nested) == 0
}

// Размер объекта.
//
// Для массива - количество элементов.
// Для объекта ключ-значение - количество ключей.
// Для скалярного значения - -1.
func (j *Nested) Length() int {
	if j.IsValue() {
		return -1
	}

	if j.IsArray() {
		return len(j.array)
	}

	return len(j.nested)
}

// Рекурсивная очистка объекта и всех вложенных.
//
// Следует учитывать, что внутри структуры используются указатели. Если структура была инициализирована
// указателями на внешние объекты, они тоже могут стать недоступны.
func (j *Nested) Clear() error {
	if j.IsValue() {
		j.isValue = false
		j.value = nil

		return nil
	}

	if j.IsArray() {
		for _, element := range j.array[:] {
			element.Clear()
			element = nil
		}

		j.isArray = false
		j.array = nil

		return nil
	}

	for k := range j.nested {
		j.nested[k].Clear()
		j.nested[k] = nil
		delete(j.nested, k)
	}

	return nil
}

// Получение указателя на вложенный объект по цепочке ключей.
//
// Должен быть указан хотя бы один ключ, иначе вернется ошибка.
// Если отсутствует один из промежуточных ключей, также вернется ошибка.
//
// Если исходный или один из промежуточных объектов является массивом или значением, функция вернет ошибку.
//
// Функцию также следует использовать для проверки наличия ключа в объекте через сравнение с nil возвращенного указателя или ошибки.
func (j *Nested) Get(keys ...string) (*Nested, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("keys list must contain at least one key")
	}

	if j.IsValue() {
		return nil, fmt.Errorf("is value")
	}

	if j.IsArray() {
		return nil, fmt.Errorf("is array")
	}

	currentKey := keys[0]

	value, ok := j.nested[currentKey]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", currentKey)
	}

	if len(keys) == 1 {
		return value, nil
	}

	res, err := value.Get(keys[1:]...)
	if err != nil {
		formatString := "%s.%s"
		if len(keys) == 2 {
			formatString = "%s: %s"
		}
		return nil, fmt.Errorf(formatString, currentKey, err.Error())
	}

	return res, err
}

// Помещение вложенного объекта по цепочке ключей.
//
// Должен быть указан хотя бы один ключ, иначе вернется ошибка.
// Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек
// будут созданы новые вложенные объекты.
//
// Если исходный или один из промежуточных объектов является массивом или значением, функция вернет ошибку.
//
// Функция принимает указатель на сохраняемый объект.
// Если в дальнейшем изменится исходный объект, изменится и вложенный.
func (j *Nested) Set(nested *Nested, keys ...string) error {
	if len(keys) == 0 {
		return fmt.Errorf("keys list must contain at least one key")
	}

	if j.IsValue() {
		return fmt.Errorf("is value")
	}

	if j.IsArray() {
		return fmt.Errorf("is array")
	}

	currentKey := keys[0]

	if j.nested == nil {
		j.nested = map[string]*Nested{}
	}

	if len(keys) == 1 {
		j.nested[currentKey] = nested
		return nil
	}

	_, ok := j.nested[currentKey]
	if !ok {
		j.nested[currentKey] = &Nested{}
	}

	err := j.nested[currentKey].Set(nested, keys[1:]...)
	if err != nil {
		formatString := "%s.%s"
		if len(keys) == 2 {
			formatString = "%s: %s"
		}
		return fmt.Errorf(formatString, currentKey, err.Error())
	}

	return err
}

// Получение скалярного значения по цепочке ключей.
//
// Все вложенные объекты до последнего в цепочке должны быть вида ключ-значение. Последний - скалярным значением.
// Можно не передавать ключи, тогда исходный объект должен быть значением.
func (j *Nested) GetValue(keys ...string) (any, error) {
	if j.IsEmpty() {
		return nil, fmt.Errorf("is empty")
	}

	if j.IsArray() {
		return nil, fmt.Errorf("is array")
	}

	if len(keys) == 0 {
		if j.IsValue() {
			return j.value, nil
		} else if j.IsNested() {
			return nil, fmt.Errorf("is nested")
		}
	}

	nested, err := j.Get(keys...)
	if err != nil {
		return nil, err
	}

	if nested.IsArray() {
		return nil, fmt.Errorf("%s: is array", strings.Join(keys[:], "."))
	}

	if nested.IsNested() {
		return nil, fmt.Errorf("%s: is nested", strings.Join(keys[:], "."))
	}

	return nested.value, nil
}

// Сохранение объекта со скалярным значением из аргумента по цепочке ключей.
//
// Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек
// будут созданы новые вложенные объекты.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
// Исходный объект может являться значением, если не передана цепочка ключей.
// В этом случае значение заменится на аргумент.
//
// Также можно не передавать цепочку ключей, если исходный объект является пустым (IsEmpty).
// В этом случае он станет объектом-значением.
func (j *Nested) SetValue(value any, keys ...string) error {
	if (j.IsEmpty() || j.IsValue()) && len(keys) == 0 {
		j.isValue = true
		j.value = value

		return nil
	}

	return j.Set(&Nested{
		isValue: true,
		value:   value,
	},
		keys...)
}

// Получение вложенного объекта вида ключ-значение (map) по цепочке ключей.
//
// Если отсутствует один из промежуточных ключей, вернется ошибка.
//
// Если один из объектов в цепочке является массивом или значением, функция вернет ошибку.
func (j *Nested) GetMap(keys ...string) (map[string]*Nested, error) {
	if j.IsValue() {
		return nil, fmt.Errorf("is value")
	}

	if j.IsArray() {
		return nil, fmt.Errorf("is array")
	}

	if len(keys) == 0 {
		return j.nested, nil
	}

	nested, err := j.Get(keys...)
	if err != nil {
		return nil, err
	}

	if nested.IsValue() {
		return nil, fmt.Errorf("%s: is value", strings.Join(keys[:], "."))
	}

	if nested.IsArray() {
		return nil, fmt.Errorf("%s: is array", strings.Join(keys[:], "."))
	}

	return nested.nested, nil
}

// Сохранение map-объекта типа map[string]*Nested по цепочке ключей.
//
// Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек
// будут созданы новые вложенные объекты.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
//
// Также можно не передавать цепочку ключей, если исходный объект является пустым (IsEmpty).
// В этом случае он станет объектом-значением.
//
// Если исходный объект непустой, удаление старых элементов не производится, и указатели на них останутся корректными.
func (j *Nested) SetMap(nested map[string]*Nested, keys ...string) error {
	if (j.IsEmpty() || j.IsNested()) && len(keys) == 0 {
		j.nested = nested

		return nil
	}

	return j.Set(&Nested{
		nested: nested,
	},
		keys...)
}

// Получение вложенного массива объектов по цепочке ключей.
//
// Все вложенные объекты до последнего в цепочке должны быть вида ключ-значение. Последний - массивом.
// Можно не передавать ключи, тогда исходный объект должен быть массивом.
func (j *Nested) GetArray(keys ...string) ([]*Nested, error) {
	if j.IsEmpty() {
		return nil, fmt.Errorf("is empty")
	}

	if j.IsValue() {
		return nil, fmt.Errorf("is value")
	}

	if len(keys) == 0 {
		if j.IsArray() {
			return j.array, nil
		} else if j.IsNested() {
			return nil, fmt.Errorf("is nested")
		}
	}

	nested, err := j.Get(keys...)
	if err != nil {
		return nil, err
	}

	if nested.IsValue() {
		return nil, fmt.Errorf("%s: is value", strings.Join(keys[:], "."))
	}

	if nested.IsNested() {
		return nil, fmt.Errorf("%s: is nested", strings.Join(keys[:], "."))
	}

	return nested.array, nil
}

// Сохранение объекта-массива из аргумента по цепочке ключей.
//
// Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек
// будут созданы новые вложенные объекты.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
// Исходный объект может являться массивом, если не передана цепочка ключей.
// В этом случае значение заменится на аргумент.
//
// Также можно не передавать цепочку ключей, если исходный объект является пустым (IsEmpty).
// В этом случае он станет объектом-массивом.
func (j *Nested) SetArray(array []*Nested, keys ...string) error {
	if (j.IsEmpty() || j.IsArray()) && len(keys) == 0 {
		j.isArray = true
		j.array = array

		return nil
	}

	return j.Set(&Nested{
		isArray: true,
		array:   array,
	},
		keys...)
}

// Удаление вложенного объекта по цепочке ключей.
//
// Должен быть передан хотя бы один ключ.
// Если будет отсутствовать один из промежуточных ключей, вернется ошибка.
//
// Все промежуточные объекты должны быть вида ключ-значение.
//
// Если последний ключ в цепочке отсутствует, функция завершится без ошибок.
func (j *Nested) Delete(keys ...string) error {
	if len(keys) == 0 {
		return fmt.Errorf("keys list must contain at least one key")
	}

	if j.IsValue() {
		return fmt.Errorf("is value")
	}

	if j.IsArray() {
		return fmt.Errorf("is array")
	}

	if len(keys) == 1 {
		j.nested[keys[0]] = nil
		delete(j.nested, keys[0])
		return nil
	}

	keysWithoutLast := keys[:len(keys)-1]

	nested, err := j.Get(keysWithoutLast...)
	if err != nil {
		return err
	}

	if nested.IsValue() {
		return fmt.Errorf("%s: is value", strings.Join(keysWithoutLast, "."))
	}

	if nested.IsArray() {
		return fmt.Errorf("%s: is array", strings.Join(keysWithoutLast, "."))
	}

	nested.nested[keys[len(keys)-1]] = nil
	delete(nested.nested, keys[len(keys)-1])

	return nil
}

// Добавление указателя на объект в массив по цепочке ключей.
//
// Если отсутствует промежуточный ключ, вернется ошибка.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
// Последний объект в цепочке должен быть массивом.
// Исходный объект может являться массивом, если не передана цепочка ключей.
func (j *Nested) ArrayAdd(element *Nested, keys ...string) error {
	if len(keys) == 0 {
		if j.IsValue() {
			return fmt.Errorf("is value")
		} else if j.IsNested() {
			return fmt.Errorf("is nested")
		}

		j.array = append(j.array, element)
		return nil
	}

	if j.IsEmpty() {
		return fmt.Errorf("is empty")
	}

	if j.IsArray() {
		return fmt.Errorf("is array")
	}

	nested, err := j.Get(keys...)
	if err != nil {
		return err
	}

	if nested.IsValue() {
		return fmt.Errorf("%s: is value", strings.Join(keys[:], "."))
	}

	if nested.IsNested() {
		return fmt.Errorf("%s: is nested", strings.Join(keys[:], "."))
	}

	nested.array = append(nested.array, element)

	return nil
}

// Добавление указателя на объект-значение по переданному аргументу в массив по цепочке ключей.
//
// Если отсутствует промежуточный ключ, вернется ошибка.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
// Последний объект в цепочке должен быть массивом.
// Исходный объект может являться массивом, если не передана цепочка ключей.
func (j *Nested) ArrayAddValue(element any, keys ...string) error {
	return j.ArrayAdd(
		&Nested{
			isValue: true,
			value:   element,
		},
		keys...,
	)
}

// Добавление указателя на объект-массив по переданному аргументу в массив по цепочке ключей.
//
// Если отсутствует промежуточный ключ, вернется ошибка.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
// Последний объект в цепочке должен быть массивом.
// Исходный объект может являться массивом, если не передана цепочка ключей.
func (j *Nested) ArrayAddArray(element []*Nested, keys ...string) error {
	return j.ArrayAdd(
		&Nested{
			isArray: true,
			array:   element,
		},
		keys...,
	)
}

// Поиск в массиве по цепочке ключей всех элементов на основе функции поиска.
//
// Если отсутствует промежуточный ключ, вернется ошибка.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
// Последний объект в цепочке должен быть массивом.
// Исходный объект может являться массивом, если не передана цепочка ключей.
//
// Если подходящих объектов в массиве нет, функция завершит работу без ошибок и вернет пустой массив.
//
// Функция не различает виды вложенных объектов.
// Для поиска только по скалярным значениям следует использовать [IsValue] в функции поиска.
// Аналогично - [IsArray] и [IsNested] для массивов и объектов ключ-значение.
//
// Пример поиска объектов-значений в объекте-массиве:
//
//	nested := Nested{}
//	nested.SetArray([]*Nested{})
//	nested.ArrayAdd(&Nested{})
//	nested.ArrayAddValue(42)
//	nested.ArrayAddValue(142)
//
//	array, _ := nested.ArrayFindAll(
//		func(element *Nested) bool {
//			return element.IsValue()
//		},
//	)
//
//	array[0].GetValue() // 42, nil
//	array[1].GetValue() // 142, nil
//
// Пример поиска объектов-массивов в объекте-массиве:
//
//	nested := Nested{}
//	nested.SetArray([]*Nested{})
//	   nested.ArrayAdd(&Nested{})
//	nested.ArrayAddValue(42)
//	nested.ArrayAddValue(142)
//
//	array, _ := nested.ArrayFindOne(
//		func(element *Nested) bool {
//			return element.IsArray()
//		},
//	) // []*Nested{}, nil
//
// Пример поиска элементов - целочисленных значений больше 100:
//
//	nested := Nested{}
//	nested.SetArray([]*Nested{})
//	nested.ArrayAdd(&Nested{})
//	nested.ArrayAddValue(42)
//	nested.ArrayAddValue(142)
//
//	array, _ := nested.ArrayFindOne(
//		func(element *Nested) bool {
//			value, err := element.GetValue()
//			if err != nil {
//				return false
//			}
//
//			v, ok := value.(int)
//			if !ok {
//				return false
//			}
//
//			return v > 100
//		},
//	)
//
//	array[0].GetValue() // 142, nil
func (j *Nested) ArrayFindAll(f func(*Nested) bool, keys ...string) ([]*Nested, error) {
	array, err := j.GetArray(keys...)
	if err != nil {
		return nil, err
	}

	found := []*Nested{}

	for _, element := range array {
		if f(element) {
			found = append(found, element)
		}
	}

	return found, nil
}

// Поиск в массиве по цепочке ключей первого подходящего элемента на основе функции поиска.
//
// Если отсутствует промежуточный ключ, вернется ошибка.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
// Последний объект в цепочке должен быть массивом.
// Исходный объект может являться массивом, если не передана цепочка ключей.
//
// Если подходящих объектов в массиве нет, функция завершит работу без ошибок и вернет нулевой указатель.
//
// Функция не различает виды вложенных объектов.
// Для поиска только по скалярным значениям следует использовать [IsValue] в функции поиска.
// Аналогично - [IsArray] и [IsNested] для массивов и объектов ключ-значение.
//
// Пример поиска первого объекта-значения в объекте-массиве:
//
//	nested := Nested{}
//	nested.SetArray([]*Nested{})
//	nested.ArrayAdd(&Nested{})
//	nested.ArrayAddValue(42)
//	nested.ArrayAddValue(142)
//
//	value, _ := nested.ArrayFindOne(
//		func(element *Nested) bool {
//			return element.IsValue()
//		},
//	)
//
//	value.GetValue() // 42, nil
//
// Пример поиска первого объекта-массива в объекте-массиве:
//
//	nested := Nested{}
//	nested.SetArray([]*Nested{})
//	nested.ArrayAdd(&Nested{})
//	nested.ArrayAddValue(42)
//	nested.ArrayAddValue(142)
//
//	value, _ := nested.ArrayFindOne(
//		func(element *Nested) bool {
//			return element.IsArray()
//		},
//	) // nil, nil
//
// Пример поиска первого элемента - положительного целочисленного значения:
//
//	nested := Nested{}
//	nested.SetArray([]*Nested{})
//	nested.ArrayAdd(&Nested{})
//	nested.ArrayAddValue(42)
//	nested.ArrayAddValue(142)
//
//	value, _ := nested.ArrayFindOne(
//		func(element *Nested) bool {
//			value, err := element.GetValue()
//			if err != nil {
//				return false
//			}
//
//			v, ok := value.(int)
//			if !ok {
//				return false
//			}
//
//			return v > 0
//		},
//	)
//
//	value.GetValue() // 42, nil
func (j *Nested) ArrayFindOne(f func(element *Nested) bool, keys ...string) (*Nested, error) {
	array, err := j.GetArray(keys...)
	if err != nil {
		return nil, err
	}

	index := slices.IndexFunc(array, f)

	if index == -1 {
		return nil, nil
	}

	return array[index], nil
}

// Удаление элементов из массива по цепочке ключей на основе функции поиска.
//
// Если отсутствует промежуточный ключ, вернется ошибка.
//
// Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.
// Последний объект в цепочке должен быть массивом.
// Исходный объект может являться массивом, если не передана цепочка ключей.
//
// Если подходящего объекта в массиве нет, функция завершит работу без ошибок.
// Будут удалены (и предварительно очищены функцией Clear() без проверки ошибок) все найденные функцией поиска элементы.
//
// Пример удаления первого объектов-значений в объекте-массиве:
//
//	nested := Nested{}
//	nested.SetArray([]*Nested{})
//	nested.ArrayAdd(&Nested{})
//	nested.ArrayAddValue(42)
//	nested.ArrayAddValue(142)
//
//	nested.ArrayDelete(
//		func(element *Nested) bool {
//			return element.IsValue()
//		},
//	)
//
//	nested.GetArray() // []*Nested{{}}, nil
func (j *Nested) ArrayDelete(f func(element *Nested) bool, keys ...string) error {
	array, err := j.GetArray(keys...)
	if err != nil {
		return err
	}

	for index := len(array) - 1; index >= 0; index-- {
		if f(array[index]) {
			array[index].Clear()
			array[index] = nil
			array = append(array[:index], array[index+1:]...)
		}
	}

	return j.SetArray(array, keys...)
}

// Создание объекта из JSON-строки.
//
// Если строка не является корректным объектом или массивом, будет возвращен объект-скалярное значение.
//
// Поддерживаются скалярные значения типов (в порядке проверки типов при конвертации) int, float64, bool, string.
//
// Примеры:
//
//	FromJSONString("[[4, 5]").GetValue()  // "[[4, 5]"
//	FromJSONString("string").GetValue()   // "string"
//	FromJSONString("42").GetValue()       // 42
//	FromJSONString("true").GetValue()     // true
func FromJSONString(nested string) *Nested {
	var kvObject map[string]any
	if err := json.Unmarshal([]byte(nested), &kvObject); err == nil {
		return FromObject(kvObject)
	}

	var arrayObject []any
	if err := json.Unmarshal([]byte(nested), &arrayObject); err == nil {
		return FromObject(arrayObject)
	}

	result := Nested{isValue: true}

	if value, err := strconv.ParseInt(nested, 10, 32); err == nil {
		result.value = int(value)
	} else if value, err := strconv.ParseFloat(nested, 32); err == nil {
		result.value = value
	} else if value, err := strconv.ParseBool(nested); err == nil {
		result.value = value
	} else {
		result.value = nested
	}

	return &result
}

// Рекурсивная функция конвертации интерфейса в Nested.
//
// Вложенные словари могут быть вида map[string]any.
// Если ключи словаря имеют другой тип, словарь сохранится как объект-значение.
//
// Если объект не является массивом []any или словарем map[string]any,
// он сохраняется как исходный тип интерфейса в объект-значение,
// кроме float64. Для него производится попытка конвертации в int.
// Это связано с тем, что функция используется в [FromJSONString],
// где в свою очередь для парсинга строки с объектом или массивом используется Unmarshal из
// пакета [ https://pkg.go.dev/encoding/json ], в котором все числа парсятся как float64.
//
// Примеры:
//
//	nested = FromObject(map[string]any{"a": 1, "b": 2})
//	nested.IsNested() // true
//	nested.IsEmpty() // true
//
//	nested = FromObject([]any{1,2,3})
//	nested.IsArray() // true
//
//	nested = FromObject(42)
//	nested.IsValue() // true
func FromObject(obj any) *Nested {
	if kvObject, ok := obj.(map[string]any); ok {
		nested := make(map[string]*Nested)

		for k := range kvObject {
			nested[k] = FromObject(kvObject[k])
		}

		return &Nested{nested: nested}
	}

	if arrayObject, ok := obj.([]any); ok {
		var array []*Nested

		for _, element := range arrayObject {
			array = append(array, FromObject(element))
		}

		return &Nested{isArray: true, array: array}
	}

	valueNested := Nested{isValue: true}

	if value, ok := obj.(float64); ok {
		if value == float64(int(value)) {
			valueNested.value = int(value)
		} else {
			valueNested.value = value
		}
	} else {
		valueNested.value = obj
	}

	return &valueNested
}

// Конвертация объекта в JSON-строку.
//
// Для словарей (map) ключи будут отсортированы в алфавитном порядке.
// Если исходный объект - строковое значение, для него будут удалены лишние обрамляющие кавычки.
//
// Пример:
//
//	nested := Nested{}
//
//	nested.SetValue("value1", "nested_object", "key1")
//	nested.SetValue("value2", "nested_object", "key2")
//	nested.SetArray([]*Nested{}, "nested_array")
//	nested.ArrayAddValue("elem1", "nested_array")
//	nested.ArrayAddValue("elem2", "nested_array")
//
//	nested.ToJSONString() // {"nested_array":["elem1","elem2"],"nested_object":{"key1":"value1","key2":"value2"}}
func (j *Nested) ToJSONString() string {
	// удаление лишних обрамляющих кавычек
	trim := func(s string) string {
		if s[len(s)-1] == '\n' {
			s = s[0 : len(s)-1]
		}

		if len(s) >= 2 {
			if s[0] == '"' && s[len(s)-1] == '"' {
				s = s[1 : len(s)-1]
			}
		}

		return s
	}

	result := ""

	// отключение экранирования спецсимволов
	jsonMarshal := func(t interface{}) ([]byte, error) {
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(t)
		return buffer.Bytes(), err
	}

	if objString, err := jsonMarshal(j.ToObject()); err == nil {
		result = trim(string(objString))
	}

	return result
}

// Конвертация объекта в объект-интерфейс.
//
// Пример:
//
//	nested := Nested{}
//
//	nested.SetValue("value1", "nested_object", "key1")
//	nested.SetValue("value2", "nested_object", "key2")
//	nested.SetArray([]*Nested{}, "nested_array")
//	nested.ArrayAddValue("elem1", "nested_array")
//	nested.ArrayAddValue("elem2", "nested_array")
//
//	// map[string]any{
//	// 	"nested_object": map[string]any{
//	// 		"key1": "value1",
//	// 		"key2": "value2",
//	// 	},
//	// 	"nested_array": []any{"elem1", "elem2"},
//	// }
//	nested.ToObject()
func (j *Nested) ToObject() any {
	if j.IsValue() {
		return j.value
	}

	if j.IsArray() {
		result := []any{}

		for _, element := range j.array {
			result = append(result, element.ToObject())
		}

		return result
	}

	result := make(map[string]any)

	for k := range j.nested {
		result[k] = j.nested[k].ToObject()
	}

	return result
}

// Функция для сравнения двух объектов Nested.
//
// Возвращает true, если объекты равны, в противном случае false.
// При сравнении учитываются не только сами значения, но и типы данных объекта.
// Также, если элементы содержатся в массиве, важен порядок их расположения, то есть в соответствующих индексах массива должны быть равные элементы.
//
// Пример:
//
//	var a Nested = Nested{isValue: true, value: int(5)}.
//	var b Nested = Nested{isValue: true, value: uint64(5)}.
//	// Equals(&a, &b) вернет false, так как переменные имеют разный тип данных.
//
//	var a Nested = Nested{
//		isArray: true,
//		array: []*Nested{
//			{
//				nested: map[string]*Nested{
//					"nested_value": {
//						isValue: true,
//						value:   "string in nested array",
//					},
//				},
//			},
//			{
//				nested: map[string]*Nested{
//					"super value": {
//						isValue: true,
//						value:   "string in nested array",
//					},
//				},
//			},
//		},
//	}
//	var b Nested = Nested{
//		isArray: true,
//		array: []*Nested{
//			{
//				nested: map[string]*Nested{
//					"super value": {
//						isValue: true,
//						value:   "string in nested array",
//					},
//				},
//			},
//			{
//				nested: map[string]*Nested{
//					"nested_value": {
//						isValue: true,
//						value:   "string in nested array",
//					},
//				},
//			},
//		},
//	}
//	// Equals(&a, &b) вернет false, так как важен порядок элементов в массиве.
func Equals(a, b *Nested) bool {
	return reflect.DeepEqual(a.ToObject(), b.ToObject())
}
