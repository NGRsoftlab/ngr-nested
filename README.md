<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# nested

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

Eсть возможность инициализации структуры из JSON\-строки и обратно с помощью [FromJSONString](<#FromJSONString>) и \[ToJSONString\]:

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

Для инициализации структуры словарем, массивом или значением\-интерфейсом и обратной конвертации в интерфейс см. [FromObject](<#FromObject>) и \[ToObject\].

## Index

- [func Equals\(a, b \*Nested\) bool](<#Equals>)
- [type Nested](<#Nested>)
  - [func FromJSONString\(nested string\) \*Nested](<#FromJSONString>)
  - [func FromObject\(obj any\) \*Nested](<#FromObject>)
  - [func \(j \*Nested\) ArrayAdd\(element \*Nested, keys ...string\) error](<#Nested.ArrayAdd>)
  - [func \(j \*Nested\) ArrayAddArray\(element \[\]\*Nested, keys ...string\) error](<#Nested.ArrayAddArray>)
  - [func \(j \*Nested\) ArrayAddValue\(element any, keys ...string\) error](<#Nested.ArrayAddValue>)
  - [func \(j \*Nested\) ArrayDelete\(f func\(element \*Nested\) bool, keys ...string\) error](<#Nested.ArrayDelete>)
  - [func \(j \*Nested\) ArrayFindAll\(f func\(\*Nested\) bool, keys ...string\) \(\[\]\*Nested, error\)](<#Nested.ArrayFindAll>)
  - [func \(j \*Nested\) ArrayFindOne\(f func\(element \*Nested\) bool, keys ...string\) \(\*Nested, error\)](<#Nested.ArrayFindOne>)
  - [func \(j \*Nested\) Clear\(\) error](<#Nested.Clear>)
  - [func \(j \*Nested\) Delete\(keys ...string\) error](<#Nested.Delete>)
  - [func \(j \*Nested\) Get\(keys ...string\) \(\*Nested, error\)](<#Nested.Get>)
  - [func \(j \*Nested\) GetArray\(keys ...string\) \(\[\]\*Nested, error\)](<#Nested.GetArray>)
  - [func \(j \*Nested\) GetMap\(keys ...string\) \(map\[string\]\*Nested, error\)](<#Nested.GetMap>)
  - [func \(j \*Nested\) GetValue\(keys ...string\) \(any, error\)](<#Nested.GetValue>)
  - [func \(j \*Nested\) IsArray\(\) bool](<#Nested.IsArray>)
  - [func \(j \*Nested\) IsEmpty\(\) bool](<#Nested.IsEmpty>)
  - [func \(j \*Nested\) IsNested\(\) bool](<#Nested.IsNested>)
  - [func \(j \*Nested\) IsValue\(\) bool](<#Nested.IsValue>)
  - [func \(j \*Nested\) Length\(\) int](<#Nested.Length>)
  - [func \(j \*Nested\) Set\(nested \*Nested, keys ...string\) error](<#Nested.Set>)
  - [func \(j \*Nested\) SetArray\(array \[\]\*Nested, keys ...string\) error](<#Nested.SetArray>)
  - [func \(j \*Nested\) SetMap\(nested map\[string\]\*Nested, keys ...string\) error](<#Nested.SetMap>)
  - [func \(j \*Nested\) SetValue\(value any, keys ...string\) error](<#Nested.SetValue>)
  - [func \(j \*Nested\) ToJSONString\(\) string](<#Nested.ToJSONString>)
  - [func \(j \*Nested\) ToObject\(\) any](<#Nested.ToObject>)


<a name="Equals"></a>
## func [Equals](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L1010>)

```go
func Equals(a, b *Nested) bool
```

Функция для сравнения двух объектов Nested.

Возвращает true, если объекты равны, в противном случае false. При сравнении учитываются не только сами значения, но и типы данных объекта. Также, если элементы содержатся в массиве, важен порядок их расположения, то есть в соответствующих индексах массива должны быть равные элементы.

Пример:

```
var a Nested = Nested{isValue: true, value: int(5)}.
var b Nested = Nested{isValue: true, value: uint64(5)}.
// Equals(&a, &b) вернет false, так как переменные имеют разный тип данных.

var a Nested = Nested{
	isArray: true,
	array: []*Nested{
		{
			nested: map[string]*Nested{
				"nested_value": {
					isValue: true,
					value:   "string in nested array",
				},
			},
		},
		{
			nested: map[string]*Nested{
				"super value": {
					isValue: true,
					value:   "string in nested array",
				},
			},
		},
	},
}
var b Nested = Nested{
	isArray: true,
	array: []*Nested{
		{
			nested: map[string]*Nested{
				"super value": {
					isValue: true,
					value:   "string in nested array",
				},
			},
		},
		{
			nested: map[string]*Nested{
				"nested_value": {
					isValue: true,
					value:   "string in nested array",
				},
			},
		},
	},
}
// Equals(&a, &b) вернет false, так как важен порядок элементов в массиве.
```

<a name="Nested"></a>
## type [Nested](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L69-L77>)

Структура для описания объекта.

Объект может быть скалярным значением любого типа \(any\), массивом указателей на объекты или объектом вида ключ\-значение, где ключ \- строка, значение \- указатель на объект.

Для каждого объекта может использоваться только один из этих вариантов поведения. По умолчанию объект имеет вид ключ\-значение. Другие варианты обозначаются выставлением флагов isValue и isArray.

Поля в объекте не экспортируются, чтобы избежать случайных конфликтов при инициализации.

```go
type Nested struct {
    isValue bool // является ли объект скалярным значением
    isArray bool // является ли объект массивом

    nested map[string]*Nested // вложенный объект вида ключ-значение
    array  []*Nested          // массив объектов

    value any // скалярное значение
}
```

<a name="FromJSONString"></a>
### func [FromJSONString](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L776>)

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

<a name="FromObject"></a>
### func [FromObject](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L825>)

```go
func FromObject(obj any) *Nested
```

Рекурсивная функция конвертации интерфейса в Nested.

Вложенные словари могут быть вида map\[string\]any. Если ключи словаря имеют другой тип, словарь сохранится как объект\-значение.

Если объект не является массивом \[\]any или словарем map\[string\]any, он сохраняется как исходный тип интерфейса в объект\-значение, кроме float64. Для него производится попытка конвертации в int. Это связано с тем, что функция используется в [FromJSONString](<#FromJSONString>), где в свою очередь для парсинга строки с объектом или массивом используется Unmarshal из пакета \[ https://pkg.go.dev/encoding/json \], в котором все числа парсятся как float64.

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

<a name="Nested.ArrayAdd"></a>
### func \(\*Nested\) [ArrayAdd](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L479>)

```go
func (j *Nested) ArrayAdd(element *Nested, keys ...string) error
```

Добавление указателя на объект в массив по цепочке ключей.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

<a name="Nested.ArrayAddArray"></a>
### func \(\*Nested\) [ArrayAddArray](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L541>)

```go
func (j *Nested) ArrayAddArray(element []*Nested, keys ...string) error
```

Добавление указателя на объект\-массив по переданному аргументу в массив по цепочке ключей.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

<a name="Nested.ArrayAddValue"></a>
### func \(\*Nested\) [ArrayAddValue](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L524>)

```go
func (j *Nested) ArrayAddValue(element any, keys ...string) error
```

Добавление указателя на объект\-значение по переданному аргументу в массив по цепочке ключей.

Если отсутствует промежуточный ключ, вернется ошибка.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Последний объект в цепочке должен быть массивом. Исходный объект может являться массивом, если не передана цепочка ключей.

<a name="Nested.ArrayDelete"></a>
### func \(\*Nested\) [ArrayDelete](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L748>)

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

<a name="Nested.ArrayFindAll"></a>
### func \(\*Nested\) [ArrayFindAll](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L621>)

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

<a name="Nested.ArrayFindOne"></a>
### func \(\*Nested\) [ArrayFindOne](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L707>)

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

<a name="Nested.Clear"></a>
### func \(\*Nested\) [Clear](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L120>)

```go
func (j *Nested) Clear() error
```

Рекурсивная очистка объекта и всех вложенных.

Следует учитывать, что внутри структуры используются указатели. Если структура была инициализирована указателями на внешние объекты, они тоже могут стать недоступны.

<a name="Nested.Delete"></a>
### func \(\*Nested\) [Delete](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L434>)

```go
func (j *Nested) Delete(keys ...string) error
```

Удаление вложенного объекта по цепочке ключей.

Должен быть передан хотя бы один ключ. Если будет отсутствовать один из промежуточных ключей, вернется ошибка.

Все промежуточные объекты должны быть вида ключ\-значение.

Если последний ключ в цепочке отсутствует, функция завершится без ошибок.

<a name="Nested.Get"></a>
### func \(\*Nested\) [Get](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L155>)

```go
func (j *Nested) Get(keys ...string) (*Nested, error)
```

Получение указателя на вложенный объект по цепочке ключей.

Должен быть указан хотя бы один ключ, иначе вернется ошибка. Если отсутствует один из промежуточных ключей, также вернется ошибка.

Если исходный или один из промежуточных объектов является массивом или значением, функция вернет ошибку.

Функцию также следует использовать для проверки наличия ключа в объекте через сравнение с nil возвращенного указателя или ошибки.

<a name="Nested.GetArray"></a>
### func \(\*Nested\) [GetArray](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L367>)

```go
func (j *Nested) GetArray(keys ...string) ([]*Nested, error)
```

Получение вложенного массива объектов по цепочке ключей.

Все вложенные объекты до последнего в цепочке должны быть вида ключ\-значение. Последний \- массивом. Можно не передавать ключи, тогда исходный объект должен быть массивом.

<a name="Nested.GetMap"></a>
### func \(\*Nested\) [GetMap](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L310>)

```go
func (j *Nested) GetMap(keys ...string) (map[string]*Nested, error)
```

Получение вложенного объекта вида ключ\-значение \(map\) по цепочке ключей.

Если отсутствует один из промежуточных ключей, вернется ошибка.

Если один из объектов в цепочке является массивом или значением, функция вернет ошибку.

<a name="Nested.GetValue"></a>
### func \(\*Nested\) [GetValue](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L246>)

```go
func (j *Nested) GetValue(keys ...string) (any, error)
```

Получение скалярного значения по цепочке ключей.

Все вложенные объекты до последнего в цепочке должны быть вида ключ\-значение. Последний \- скалярным значением. Можно не передавать ключи, тогда исходный объект должен быть значением.

<a name="Nested.IsArray"></a>
### func \(\*Nested\) [IsArray](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L80>)

```go
func (j *Nested) IsArray() bool
```

Проверка, является ли объект массивом.

<a name="Nested.IsEmpty"></a>
### func \(\*Nested\) [IsEmpty](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L95>)

```go
func (j *Nested) IsEmpty() bool
```

Проверка, что объект пустой и имеет вид ключ\-значение.

<a name="Nested.IsNested"></a>
### func \(\*Nested\) [IsNested](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L90>)

```go
func (j *Nested) IsNested() bool
```

Проверка, что объект имеет тип ключ\-значение.

<a name="Nested.IsValue"></a>
### func \(\*Nested\) [IsValue](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L85>)

```go
func (j *Nested) IsValue() bool
```

Проверка, является ли объект скалярным значением.

<a name="Nested.Length"></a>
### func \(\*Nested\) [Length](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L104>)

```go
func (j *Nested) Length() int
```

Размер объекта.

Для массива \- количество элементов. Для объекта ключ\-значение \- количество ключей. Для скалярного значения \- \-1.

<a name="Nested.Set"></a>
### func \(\*Nested\) [Set](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L201>)

```go
func (j *Nested) Set(nested *Nested, keys ...string) error
```

Помещение вложенного объекта по цепочке ключей.

Должен быть указан хотя бы один ключ, иначе вернется ошибка. Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек будут созданы новые вложенные объекты.

Если исходный или один из промежуточных объектов является массивом или значением, функция вернет ошибку.

Функция принимает указатель на сохраняемый объект. Если в дальнейшем изменится исходный объект, изменится и вложенный.

<a name="Nested.SetArray"></a>
### func \(\*Nested\) [SetArray](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L411>)

```go
func (j *Nested) SetArray(array []*Nested, keys ...string) error
```

Сохранение объекта\-массива из аргумента по цепочке ключей.

Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек будут созданы новые вложенные объекты.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Исходный объект может являться массивом, если не передана цепочка ключей. В этом случае значение заменится на аргумент.

Также можно не передавать цепочку ключей, если исходный объект является пустым \(IsEmpty\). В этом случае он станет объектом\-массивом.

<a name="Nested.SetMap"></a>
### func \(\*Nested\) [SetMap](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L350>)

```go
func (j *Nested) SetMap(nested map[string]*Nested, keys ...string) error
```

Сохранение map\-объекта типа map\[string\]\*Nested по цепочке ключей.

Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек будут созданы новые вложенные объекты.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку.

Также можно не передавать цепочку ключей, если исходный объект является пустым \(IsEmpty\). В этом случае он станет объектом\-значением.

Если исходный объект непустой, удаление старых элементов не производится, и указатели на них останутся корректными.

<a name="Nested.SetValue"></a>
### func \(\*Nested\) [SetValue](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L290>)

```go
func (j *Nested) SetValue(value any, keys ...string) error
```

Сохранение объекта со скалярным значением из аргумента по цепочке ключей.

Если отсутствует промежуточный ключ, для него и для всех остальных ключей в цепочек будут созданы новые вложенные объекты.

Если один из вложенных объектов в цепочке является массивом или значением, функция вернет ошибку. Исходный объект может являться значением, если не передана цепочка ключей. В этом случае значение заменится на аргумент.

Также можно не передавать цепочку ключей, если исходный объект является пустым \(IsEmpty\). В этом случае он станет объектом\-значением.

<a name="Nested.ToJSONString"></a>
### func \(\*Nested\) [ToJSONString](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L877>)

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

<a name="Nested.ToObject"></a>
### func \(\*Nested\) [ToObject](<https://github.com/smirnoffkin/ngr-nested/blob/main/nested.go#L931>)

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

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
