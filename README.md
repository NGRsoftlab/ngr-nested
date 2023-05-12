# ngr-nested

```go
import "github.com/NGRsoftlab/ngr-nested"
```

Вложенные объекты в Go.

Используются как альтернатива словарям с динамической типизацией для реализации работы с JSON\-объектами.

Не являются конкурентно\-безопасными, для добавления этой возможности надо делать обертку.

Пример использования для создания и чтения объекта:

```
// {
// 	 "nested_object": {
// 	   "key1": "value1",
//	   "key2": "value2"
//	 },
//	 "nested_array": ["elem1", "elem2"]
// }

nested := Nested{}

nested.SetValue("value1", "nested_object", "key1")
nested.SetValue("value2", "nested_object", "key2")
nested.SetArray([]*Nested{}, "nested_array")
nested.ArrayAddValue("elem1", "nested_array")
nested.ArrayAddValue("elem2", "nested_array")

nested.GetValue("nested_object", "key1")
nested.GetValue("nested_object", "key2")
array, _ := nested.GetArray("nested_array")
for _, elem := range array {
	elem.GetValue()
}
```

Eсть возможность инициализации структуры из JSON\-строки и обратно с помощью \[FromJSONString\] и \[ToJSONString\]:

```
nested := FromJSONString(`{"nested_object": {"key1": "value1", "key2": "value2"}, "nested_array": ["elem1", "elem2"]}`)

nested.GetValue("nested_object", "key1")
nested.GetValue("nested_object", "key2")
array, _ := nested.GetArray("nested_array")
for _, elem := range array {
	elem.GetValue()
}

nested.ToJSONString() // {"nested_object": {"key1": "value1", "key2": "value2"}, "nested_array": ["elem1", "elem2"]}
```

Для инициализации структуры словарем, массивом или значением\-интерфейсом и обратной конвертации в интерфейс см. \[FromObject\] и \[ToObject\].

## Index

- [type Nested](<#type-nested>)
  - [func FromObject(obj any) *Nested](<#func-fromobject>)
  - [func FromJSONString(nested string) *Nested](<#func-FromJSONString>)
  - [func (j *Nested) ArrayAdd(element *Nested, keys ...string) error](<#func-nested-arrayadd>)
  - [func (j *Nested) ArrayAddArray(element []*Nested, keys ...string) error](<#func-nested-arrayaddarray>)
  - [func (j *Nested) ArrayAddValue(element any, keys ...string) error](<#func-nested-arrayaddvalue>)
  - [func (j *Nested) ArrayDelete(f func(element *Nested) bool, keys ...string) error](<#func-nested-arraydelete>)
  - [func (j *Nested) ArrayFindAll(f func(*Nested) bool, keys ...string) ([]*Nested, error)](<#func-nested-arrayfindall>)
  - [func (j *Nested) ArrayFindOne(f func(element *Nested) bool, keys ...string) (*Nested, error)](<#func-nested-arrayfindone>)
  - [func (j *Nested) Clear() error](<#func-nested-clear>)
  - [func (j *Nested) Delete(keys ...string) error](<#func-nested-delete>)
  - [func (j *Nested) Get(keys ...string) (*Nested, error)](<#func-nested-get>)
  - [func (j *Nested) GetArray(keys ...string) ([]*Nested, error)](<#func-nested-getarray>)
  - [func (j *Nested) GetMap(keys ...string) (map[string]*Nested, error)](<#func-nested-getmap>)
  - [func (j *Nested) GetValue(keys ...string) (any, error)](<#func-nested-getvalue>)
  - [func (j *Nested) IsArray() bool](<#func-nested-isarray>)
  - [func (j *Nested) IsEmpty() bool](<#func-nested-isempty>)
  - [func (j *Nested) IsNested() bool](<#func-nested-isnested>)
  - [func (j *Nested) IsValue() bool](<#func-nested-isvalue>)
  - [func (j *Nested) Length() int](<#func-nested-length>)
  - [func (j *Nested) Set(nested *Nested, keys ...string) error](<#func-nested-set>)
  - [func (j *Nested) SetArray(array []*Nested, keys ...string) error](<#func-nested-setarray>)
  - [func (j *Nested) SetMap(nested map[string]*Nested, keys ...string) error](<#func-nested-setmap>)
  - [func (j *Nested) SetValue(value any, keys ...string) error](<#func-nested-setvalue>)
  - [func (j *Nested) ToObject() any](<#func-nested-toobject>)
  - [func (j *Nested) ToJSONString() string](<#func-nested-ToJSONString>)


## type [Nested](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L67-L75>)

Структура для описания объекта.

Объект может быть скалярным значением любого типа \(any\), массивом указателей на объекты или объектом вида ключ\-значение, где ключ \- строка, значение \- указатель на объект.

Для каждого объекта может использоваться только один из этих вариантов поведения. По умолчанию объект имеет вид ключ\-значение. Другие варианты обозначаются выставлением флагов isValue и isArray.

Поля в объекте не экспортируются, чтобы избежать случайных конфликтов при инициализации.

```go
type Nested struct {
    // contains filtered or unexported fields
}
```

### func [FromObject](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L821>)

```go
func FromObject(obj any) *Nested
```

Рекурсивная функция конвертации интерфейса в Nested.

Вложенные словари могут быть вида map\[string\]any. Если ключи словаря имеют другой тип, словарь сохранится как объект\-значение.

Если объект не является массивом \[\]any или словарем map\[string\]any, он сохраняется как исходный тип интерфейса в объект\-значение, кроме float64. Для него производится попытка конвертации в int. Это связано с тем, что функция используется в \[FromJSONString\], где в свою очередь для парсинга строки с объектом или массивом используется Unmarshal из пакета \[ https://pkg.go.dev/encoding/json \], в котором все числа парсятся как float64.

Примеры:

```
nested = FromObject(map[string]any{"a": 1, "b": 2})
nested.IsNested() // true
nested.IsEmpty() // true

nested = FromObject([]any{1,2,3})
nested.IsArray() // true

nested = FromObject(42)
nested.IsValue() // true
```

### func [FromJSONString](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L772>)

```go
func FromJSONString(nested string) *Nested
```

Создание объекта из JSON\-строки.

Если строка не является корректным объектом или массивом, будет возвращен объект\-скалярное значение.

Поддерживаются скалярные значения типов \(в порядке проверки типов при конвертации\) int, float64, bool, string.

Примеры:

```
FromJSONString("[[4, 5]").GetValue()  // "[[4, 5]"
FromJSONString("string").GetValue()   // "string"
FromJSONString("42").GetValue()       // 42
FromJSONString("true").GetValue()     // true
```

### func \(\*Nested\) [ArrayAdd](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L475>)

```go
func (j *Nested) ArrayAdd(element *Nested, keys ...string) error
```

Добавление указателя на объект в массив по цепочке ключей.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

### func \(\*Nested\) [ArrayAddArray](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L537>)

```go
func (j *Nested) ArrayAddArray(element []*Nested, keys ...string) error
```

Добавление указателя на объект\-массив по переданному аргументу в массив по цепочке ключей.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

### func \(\*Nested\) [ArrayAddValue](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L520>)

```go
func (j *Nested) ArrayAddValue(element any, keys ...string) error
```

Добавление указателя на объект\-значение по переданному аргументу в массив по цепочке ключей.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

### func \(\*Nested\) [ArrayDelete](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L744>)

```go
func (j *Nested) ArrayDelete(f func(element *Nested) bool, keys ...string) error
```

Удаление элементов из массива по цепочке ключей на основе функции поиска.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

Если подходящего объекта в массиве нет, функция завершит работу без ошибок. Будут удалены \(и предварительно очищены функцией Clear\(\) без проверки ошибок\) все найденные функцией поиска элементы.

Пример удаления первого объектов\-значений в объекте\-массиве:

```
nested := Nested{}
nested.SetArray([]*Nested{})
nested.ArrayAdd(&Nested{})
nested.ArrayAddValue(42)
nested.ArrayAddValue(142)

nested.ArrayDelete(
	func(element *Nested) bool {
		return element.IsValue()
	},
)

nested.GetArray() // []*Nested{{}}, nil
```

### func \(\*Nested\) [ArrayFindAll](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L617>)

```go
func (j *Nested) ArrayFindAll(f func(*Nested) bool, keys ...string) ([]*Nested, error)
```

Поиск в массиве по цепочке ключей всех элементов на основе функции поиска.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

Если подходящих объектов в массиве нет, функция завершит работу без ошибок и вернет пустой массив.

Функция не различает виды вложенных объектов. Для поиска только по скалярным значениям следует использовать \[IsValue\] в функции поиска. Аналогично \- \[IsArray\] и \[IsNested\] для массивов и объектов ключ\-значение.

Пример поиска объектов\-значений в объекте\-массиве:

```
nested := Nested{}
nested.SetArray([]*Nested{})
nested.ArrayAdd(&Nested{})
nested.ArrayAddValue(42)
nested.ArrayAddValue(142)

array, _ := nested.ArrayFindAll(
	func(element *Nested) bool {
		return element.IsValue()
	},
)

array[0].GetValue() // 42, nil
array[1].GetValue() // 142, nil
```

Пример поиска объектов\-массивов в объекте\-массиве:

```
nested := Nested{}
nested.SetArray([]*Nested{})
   nested.ArrayAdd(&Nested{})
nested.ArrayAddValue(42)
nested.ArrayAddValue(142)

array, _ := nested.ArrayFindOne(
	func(element *Nested) bool {
		return element.IsArray()
	},
) // []*Nested{}, nil
```

Пример поиска элементов \- целочисленных значений больше 100:

```
nested := Nested{}
nested.SetArray([]*Nested{})
nested.ArrayAdd(&Nested{})
nested.ArrayAddValue(42)
nested.ArrayAddValue(142)

array, _ := nested.ArrayFindOne(
	func(element *Nested) bool {
		value, err := element.GetValue()
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v > 100
	},
)

array[0].GetValue() // 142, nil
```

### func \(\*Nested\) [ArrayFindOne](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L703>)

```go
func (j *Nested) ArrayFindOne(f func(element *Nested) bool, keys ...string) (*Nested, error)
```

Поиск в массиве по цепочке ключей первого подходящего элемента на основе функции поиска.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

Если подходящих объектов в массиве нет, функция завершит работу без ошибок и вернет нулевой указатель.

Функция не различает виды вложенных объектов. Для поиска только по скалярным значениям следует использовать \[IsValue\] в функции поиска. Аналогично \- \[IsArray\] и \[IsNested\] для массивов и объектов ключ\-значение.

Пример поиска первого объекта\-значения в объекте\-массиве:

```
nested := Nested{}
nested.SetArray([]*Nested{})
nested.ArrayAdd(&Nested{})
nested.ArrayAddValue(42)
nested.ArrayAddValue(142)

value, _ := nested.ArrayFindOne(
	func(element *Nested) bool {
		return element.IsValue()
	},
)

value.GetValue() // 42, nil
```

Пример поиска первого объекта\-массива в объекте\-массиве:

```
nested := Nested{}
nested.SetArray([]*Nested{})
nested.ArrayAdd(&Nested{})
nested.ArrayAddValue(42)
nested.ArrayAddValue(142)

value, _ := nested.ArrayFindOne(
	func(element *Nested) bool {
		return element.IsArray()
	},
) // nil, nil
```

Пример поиска первого элемента \- положительного целочисленного значения:

```
nested := Nested{}
nested.SetArray([]*Nested{})
nested.ArrayAdd(&Nested{})
nested.ArrayAddValue(42)
nested.ArrayAddValue(142)

value, _ := nested.ArrayFindOne(
	func(element *Nested) bool {
		value, err := element.GetValue()
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v > 0
	},
)

value.GetValue() // 42, nil
```

### func \(\*Nested\) [Clear](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L118>)

```go
func (j *Nested) Clear() error
```

Рекурсивная очистка объекта и всех вложенных.

Следует учитывать, что внутри структуры используются указатели. Если структура была инициализирована указателями на внешние объекты, они тоже могут стать недоступны.

### func \(\*Nested\) [Delete](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L430>)

```go
func (j *Nested) Delete(keys ...string) error
```

Удаление вложенного объекта по цепочке ключей.

Должен быть передан хотя бы один ключ. Если будет отсутствовать один из промежуточных ключей, вернется ошибка.

Все промежуточные объекты должны быть вида ключ\-значение.

Если последний ключ в цепочке отсутствует, функция завершится без ошибок.

Функцию также следует использовать для проверки наличия ключа в объекте через сравнение с nil возвращенного указателя или ошибки.

### func \(\*Nested\) [Get](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L151>)

```go
func (j *Nested) Get(keys ...string) (*Nested, error)
```

Получение указателя на вложенный объект по цепочке ключей.

Должен быть указан хотя бы один ключ, иначе вернется ошибка. Если отсутствует один из промежуточных ключей, также вернется ошибка.

Если исходный или один из промежуточных объектов является массивом или значением, функция вернет ошибку.

### func \(\*Nested\) [GetArray](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L363>)

```go
func (j *Nested) GetArray(keys ...string) ([]*Nested, error)
```

Получение вложенного массива объектов по цепочке ключей.

Все вложенные объекты до последнего в цепочке должны быть вида ключ\-значение. Последний \- массивом. Можно не передавать ключи, тогда исходный объект должен быть массивом.

### func \(\*Nested\) [GetMap](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L306>)

```go
func (j *Nested) GetMap(keys ...string) (map[string]*Nested, error)
```

Получение вложенного объекта вида ключ\-значение \(map\) по цепочке ключей.

Если отсутствует один из промежуточных ключей, вернется ошибка.

Если один из объектов в цепочке является массивом или значением, функция вернет ошибку.

### func \(\*Nested\) [GetValue](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L242>)

```go
func (j *Nested) GetValue(keys ...string) (any, error)
```

Получение скалярного значения по цепочке ключей.

Все вложенные объекты до последнего в цепочке должны быть вида ключ\-значение. Последний \- скалярным значением. Можно не передавать ключи, тогда исходный объект должен быть значением.

### func \(\*Nested\) [IsArray](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L78>)

```go
func (j *Nested) IsArray() bool
```

Проверка, является ли объект массивом.

### func \(\*Nested\) [IsEmpty](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L93>)

```go
func (j *Nested) IsEmpty() bool
```

Проверка, что объект пустой и имеет вид ключ\-значение.

### func \(\*Nested\) [IsNested](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L88>)

```go
func (j *Nested) IsNested() bool
```

Проверка, что объект имеет тип ключ\-значение.

### func \(\*Nested\) [IsValue](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L83>)

```go
func (j *Nested) IsValue() bool
```

Проверка, является ли объект скалярным значением.

### func \(\*Nested\) [Length](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L102>)

```go
func (j *Nested) Length() int
```

Размер объекта.

Для массива \- количество элементов. Для объекта ключ\-значение \- количество ключей. Для скалярного значения \- \-1.

### func \(\*Nested\) [Set](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L197>)

```go
func (j *Nested) Set(nested *Nested, keys ...string) error
```

Помещение вложенного объекта по цепочке ключей.

Должен быть указан хотя бы один ключ, иначе вернется ошибка. Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек будут созданы новые вложенные объекты.

Если исходный или один из промежуточных объектов является массивом или значением, функция вернет ошибку.

Функция принимает указатель на сохраняемый объект. Если в дальнейшем изменится исходный объект, изменится и вложенный.

### func \(\*Nested\) [SetArray](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L407>)

```go
func (j *Nested) SetArray(array []*Nested, keys ...string) error
```

Сохранение объекта\-массива из аргумента по цепочке ключей.

Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек будут созданы новые вложенные объекты.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Исходный объект может являться массивом, если не передана цепочка ключей. В этом случае значение заменится на аргумент.

Также можно не передавать цепочку ключей, если исходный объект является пустым \(IsEmpty\). В этом случае он станет объектом\-массивом.

### func \(\*Nested\) [SetMap](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L346>)

```go
func (j *Nested) SetMap(nested map[string]*Nested, keys ...string) error
```

Сохранение map\-объекта типа map\[string\]\*Nested по цепочке ключей.

Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек будут созданы новые вложенные объекты.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.

Также можно не передавать цепочку ключей, если исходный объект является пустым \(IsEmpty\). В этом случае он станет объектом\-значением.

Если исходный объект непустой, удаление старых элементов не производится, и указатели на них останутся корректными.

### func \(\*Nested\) [SetValue](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L286>)

```go
func (j *Nested) SetValue(value any, keys ...string) error
```

Сохранение объекта со скалярным значением из аргумента по цепочке ключей.

Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек будут созданы новые вложенные объекты.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Исходный объект может являться значением, если не передана цепочка ключей. В этом случае значение заменится на аргумент.

Также можно не передавать цепочку ключей, если исходный объект является пустым \(IsEmpty\). В этом случае он станет объектом\-значением.

### func \(\*Nested\) [ToObject](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L913>)

```go
func (j *Nested) ToObject() any
```

Конвертация объекта в объект\-интерфейс.

Пример:

```
nested := Nested{}

nested.SetValue("value1", "nested_object", "key1")
nested.SetValue("value2", "nested_object", "key2")
nested.SetArray([]*Nested{}, "nested_array")
nested.ArrayAddValue("elem1", "nested_array")
nested.ArrayAddValue("elem2", "nested_array")

// map[string]any{
// 	"nested_object": map[string]any{
// 		"key1": "value1",
// 		"key2": "value2",
// 	},
// 	"nested_array": []any{"elem1", "elem2"},
// }
nested.ToObject()
```

### func \(\*Nested\) [ToJSONString](<https://github.com/NGRsoftlab/ngr-nested/blob/master/nested.go#L873>)

```go
func (j *Nested) ToJSONString() string
```

Конвертация объекта в JSON\-строку.

Для словарей \(map\) ключи будут отсортированы в алфавитном порядке. Если исходный объект \- строковое значение, для него будут удалены лишние обрамляющие кавычки.

Пример:

```
nested := Nested{}

nested.SetValue("value1", "nested_object", "key1")
nested.SetValue("value2", "nested_object", "key2")
nested.SetArray([]*Nested{}, "nested_array")
nested.ArrayAddValue("elem1", "nested_array")
nested.ArrayAddValue("elem2", "nested_array")

nested.ToJSONString() // {"nested_array":["elem1","elem2"],"nested_object":{"key1":"value1","key2":"value2"}}
```
