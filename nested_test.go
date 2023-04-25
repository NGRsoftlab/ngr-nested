package nested

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тестовый объект.
func testNested() Nested {
	return Nested{
		nested: map[string]*Nested{
			"value": {
				isValue: true,
				value:   42,
			},
			"array": {
				isArray: true,
				array: []*Nested{
					{
						isValue: true,
						value:   142,
					},
					{
						isValue: true,
						value:   "string in array",
					},
				},
			},
			"nested": {
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   "string",
					},
					"nested": {
						nested: map[string]*Nested{
							"value": {
								isValue: true,
								value:   "string in nested",
							},
						},
					},
					"array": {
						isArray: true,
						array: []*Nested{
							{
								nested: map[string]*Nested{
									"value": {
										isValue: true,
										value:   242,
									},
								},
							},
							{
								isArray: true,
								array: []*Nested{
									{
										nested: map[string]*Nested{
											"value": {
												isValue: true,
												value:   "string in nested array",
											},
										},
									},
								},
							},
							{
								nested: map[string]*Nested{
									"value": {
										isValue: true,
										value:   1242,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func Test_Checks(t *testing.T) {
	nested := Nested{
		isValue: true,
		value:   42,
	}

	assert.True(t, nested.IsValue())
	assert.False(t, nested.IsArray())
	assert.False(t, nested.IsNested())
	assert.False(t, nested.IsEmpty())

	nested = Nested{
		isArray: true,
		array: []*Nested{
			{
				isValue: true,
				value:   142,
			},
		},
	}

	assert.False(t, nested.IsValue())
	assert.True(t, nested.IsArray())
	assert.False(t, nested.IsNested())
	assert.False(t, nested.IsEmpty())

	nested = Nested{}

	assert.False(t, nested.IsValue())
	assert.False(t, nested.IsArray())
	assert.True(t, nested.IsNested())
	assert.True(t, nested.IsEmpty())

	nested = Nested{
		nested: map[string]*Nested{
			"value": {
				isValue: true,
				value:   42,
			},
		},
	}

	assert.False(t, nested.IsValue())
	assert.False(t, nested.IsArray())
	assert.True(t, nested.IsNested())
	assert.False(t, nested.IsEmpty())
}

func Test_Length(t *testing.T) {
	nested := Nested{
		isValue: true,
		value:   42,
	}

	assert.Equal(t, -1, nested.Length())

	nested = Nested{
		isArray: true,
		array: []*Nested{
			{
				isValue: true,
				value:   142,
			},
			{
				isValue: true,
				value:   "string in array",
			},
		},
	}

	assert.Equal(t, 2, nested.Length())

	nested = testNested()

	assert.Equal(t, 3, nested.Length())

	nested = Nested{}

	assert.Equal(t, 0, nested.Length())
}

func Test_Clear(t *testing.T) {
	nested := testNested()
	array := Nested{
		isArray: true,
		array: []*Nested{
			{
				isValue: true,
				value:   42,
			},
			{
				isValue: true,
				value:   "string",
			},
		},
	}
	value := Nested{
		isValue: true,
		value:   42,
	}
	clear := Nested{}

	for _, obj := range []Nested{nested, array, value, clear} {
		if assert.Nil(t, obj.Clear()) {
			assert.True(t, obj.IsEmpty())
		}
	}

	newNested := Nested{
		isValue: true,
		value:   142,
	}

	clear.Set(&newNested, "value")
	clear.Clear()

	assert.True(t, clear.IsEmpty())
	assert.True(t, newNested.IsEmpty())
}

func Test_GetValue(t *testing.T) {
	nested := testNested()

	_, err := nested.GetValue("somekey")
	assert.EqualError(t, err, "key 'somekey' not found")

	_, err = nested.GetValue("nested", "somekey")
	assert.EqualError(t, err, "nested: key 'somekey' not found")

	_, err = nested.GetValue("array")
	assert.EqualError(t, err, "array: is array")

	_, err = nested.GetValue()
	assert.EqualError(t, err, "is nested")

	value, err := nested.GetValue("value")
	if assert.Nil(t, err) {
		assert.Equal(t, 42, value.(int))
	}

	value, err = nested.GetValue("nested", "value")
	if assert.Nil(t, err) {
		assert.Equal(t, "string", value.(string))
	}

	value, err = nested.GetValue("nested", "nested", "value")
	if assert.Nil(t, err) {
		assert.Equal(t, "string in nested", value.(string))
	}

	_, err = nested.GetValue("nested", "nested", "somekey")
	assert.EqualError(t, err, "nested.nested: key 'somekey' not found")

	_, err = nested.GetValue("nested", "nested")
	assert.EqualError(t, err, "nested.nested: is nested")

	nested = Nested{
		isValue: true,
		value:   42,
	}

	value, err = nested.GetValue()
	if assert.Nil(t, err) {
		assert.Equal(t, 42, value.(int))
	}

	_, err = nested.GetValue("somekey")
	assert.EqualError(t, err, "is value")

	nested = Nested{
		isArray: true,
		array:   []*Nested{},
	}

	_, err = nested.GetValue("somekey")
	assert.EqualError(t, err, "is array")

	nested = Nested{}

	_, err = nested.GetValue()
	assert.EqualError(t, err, "is empty")
}

func Test_Get(t *testing.T) {
	nested := testNested()

	_, err := nested.Get("somekey")
	assert.EqualError(t, err, "key 'somekey' not found")

	_, err = nested.Get("nested", "somekey")
	assert.EqualError(t, err, "nested: key 'somekey' not found")

	_, err = nested.Get("nested", "nested", "somekey")
	assert.EqualError(t, err, "nested.nested: key 'somekey' not found")

	_, err = nested.Get()
	assert.EqualError(t, err, "keys list must contain at least one key")

	_, err = nested.Get("array", "somekey")
	assert.EqualError(t, err, "array: is array")

	_, err = nested.Get("nested", "nested", "value", "somekey")
	assert.EqualError(t, err, "nested.nested.value: is value")

	value, err := nested.Get("nested", "nested", "value")
	if assert.Nil(t, err) {
		assert.Equal(t, "string in nested", value.value.(string))
	}

	value, err = nested.Get("value")
	if assert.Nil(t, err) {
		assert.True(t, value.IsValue())
		assert.Equal(t, 42, value.value.(int))
	}

	value, err = nested.Get("array")
	if assert.Nil(t, err) {
		assert.True(t, value.IsArray())
		assert.Len(t, value.array, 2)
	}

	value, err = nested.Get("nested")
	if assert.Nil(t, err) {
		assert.True(t, value.IsNested())
		assert.Len(t, value.nested, 3)
	}

	value, err = nested.Get("nested", "nested")
	if assert.Nil(t, err) {
		assert.True(t, value.IsNested())
		assert.Len(t, value.nested, 1)
	}

	nested = Nested{
		isArray: true,
		array:   []*Nested{},
	}

	_, err = nested.Get("nested")
	assert.EqualError(t, err, "is array")

	nested = Nested{
		isArray: true,
		array:   []*Nested{},
	}

	_, err = nested.Get("nested")
	assert.EqualError(t, err, "is array")

	nested = Nested{
		isValue: true,
		value:   42,
	}

	_, err = nested.Get("nested")
	assert.EqualError(t, err, "is value")
}

func Test_Set(t *testing.T) {
	nested := testNested()

	newNested := &Nested{
		nested: map[string]*Nested{
			"value": {
				isValue: true,
				value:   "string in new nested",
			},
			"nested": {
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   1042,
					},
				},
			},
		},
	}

	err := nested.Set(newNested)
	assert.EqualError(t, err, "keys list must contain at least one key")

	err = nested.Set(newNested, "array", "somekey")
	assert.EqualError(t, err, "array: is array")

	err = nested.Set(newNested, "nested", "nested", "value", "somekey")
	assert.EqualError(t, err, "nested.nested.value: is value")

	err = nested.Set(newNested, "somekey")
	if assert.Nil(t, err) {
		value, err := nested.Get("somekey", "nested", "value")
		if assert.Nil(t, err) {
			assert.Equal(t, 1042, value.value.(int))
		}

		value, err = nested.Get("somekey", "value")
		if assert.Nil(t, err) {
			assert.Equal(t, "string in new nested", value.value.(string))
		}
	}

	err = nested.Set(newNested, "nested", "nested", "somekey")
	if assert.Nil(t, err) {
		value, err := nested.Get("nested", "nested", "somekey", "nested", "value")
		if assert.Nil(t, err) {
			assert.Equal(t, 1042, value.value.(int))
		}

		value, err = nested.Get("nested", "nested", "somekey", "value")
		if assert.Nil(t, err) {
			assert.Equal(t, "string in new nested", value.value.(string))
		}
	}

	clear := Nested{}
	err = clear.Set(newNested, "nested", "nested", "somekey")
	if assert.Nil(t, err) {
		value, err := nested.Get("nested", "nested", "somekey", "nested", "value")
		if assert.Nil(t, err) {
			assert.Equal(t, 1042, value.value.(int))
		}

		value, err = nested.Get("nested", "nested", "somekey", "value")
		if assert.Nil(t, err) {
			assert.Equal(t, "string in new nested", value.value.(string))
		}
	}
}

func Test_SetValue(t *testing.T) {
	nested := testNested()

	err := nested.SetValue(42, "array", "somekey")
	assert.EqualError(t, err, "array: is array")

	err = nested.SetValue(42, "nested", "nested", "value", "somekey")
	assert.EqualError(t, err, "nested.nested.value: is value")

	err = nested.SetValue(1042, "somekey")
	if assert.Nil(t, err) {
		value, err := nested.GetValue("somekey")
		if assert.Nil(t, err) {
			assert.Equal(t, 1042, value.(int))
		}
	}

	err = nested.SetValue("string in new nested", "nested", "nested", "somekey")
	if assert.Nil(t, err) {
		value, err := nested.GetValue("nested", "nested", "somekey")
		if assert.Nil(t, err) {
			assert.Equal(t, "string in new nested", value.(string))
		}
	}

	nested = Nested{}

	err = nested.SetValue("value")
	if assert.Nil(t, err) {
		value, err := nested.GetValue()
		if assert.Nil(t, err) {
			assert.Equal(t, "value", value.(string))
		}
	}
}

func Test_GetMap(t *testing.T) {
	nested := testNested()

	_, err := nested.GetMap("somekey")
	assert.EqualError(t, err, "key 'somekey' not found")

	_, err = nested.GetMap("nested", "somekey")
	assert.EqualError(t, err, "nested: key 'somekey' not found")

	_, err = nested.GetMap("value")
	assert.EqualError(t, err, "value: is value")

	_, err = nested.GetMap("array")
	assert.EqualError(t, err, "array: is array")

	obj, err := nested.GetMap()
	if assert.Nil(t, err) {
		assert.Equal(t, nested.nested, obj)
	}

	obj, err = nested.GetMap("nested")
	if assert.Nil(t, err) {
		assert.Equal(t, nested.nested["nested"].nested, obj)
	}

	_, err = nested.GetMap("nested", "array")
	assert.EqualError(t, err, "nested.array: is array")

	_, err = nested.GetMap("nested", "nested", "somekey")
	assert.EqualError(t, err, "nested.nested: key 'somekey' not found")

	obj, err = nested.GetMap("nested", "nested")
	if assert.Nil(t, err) {
		assert.Equal(t, nested.nested["nested"].nested["nested"].nested, obj)
	}

	nested = Nested{
		isValue: true,
		value:   42,
	}

	_, err = nested.GetMap()
	assert.EqualError(t, err, "is value")

	nested = Nested{
		isArray: true,
		array: []*Nested{
			{
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   "string in nested array",
					},
				},
			},
		},
	}

	_, err = nested.GetMap()
	assert.EqualError(t, err, "is array")

	nested = Nested{}

	obj, err = nested.GetMap()
	if assert.Nil(t, err) {
		for range obj {
			assert.Fail(t, "must be empty")
		}
	}
}

func Test_SetMap(t *testing.T) {
	nested := testNested()

	obj := map[string]*Nested{
		"one": {
			isValue: true,
			value:   42,
		},
		"two": {
			isValue: true,
			value:   43,
		},
	}

	err := nested.SetMap(obj, "array", "somekey")
	assert.EqualError(t, err, "array: is array")

	err = nested.SetMap(obj, "nested", "nested", "value", "somekey")
	assert.EqualError(t, err, "nested.nested.value: is value")

	err = nested.SetMap(obj, "nested", "somekey")
	if assert.Nil(t, err) {
		assert.Equal(t, obj, nested.nested["nested"].nested["somekey"].nested)
	}

	err = nested.SetMap(obj, "nested", "nested", "somekey")
	if assert.Nil(t, err) {
		assert.Equal(t, obj, nested.nested["nested"].nested["nested"].nested["somekey"].nested)
	}

	nested = Nested{}

	err = nested.SetMap(obj)
	if assert.Nil(t, err) {
		assert.Equal(t, obj, nested.nested)
	}
}

func Test_GetArray(t *testing.T) {
	nested := testNested()

	_, err := nested.GetArray("somekey")
	assert.EqualError(t, err, "key 'somekey' not found")

	_, err = nested.GetArray("nested", "somekey")
	assert.EqualError(t, err, "nested: key 'somekey' not found")

	_, err = nested.GetArray("value")
	assert.EqualError(t, err, "value: is value")

	_, err = nested.GetArray()
	assert.EqualError(t, err, "is nested")

	array, err := nested.GetArray("array")
	if assert.Nil(t, err) {
		assert.ElementsMatch(t, []*Nested{
			{
				isValue: true,
				value:   142,
			},
			{
				isValue: true,
				value:   "string in array",
			},
		}, array)
	}

	array, err = nested.GetArray("nested", "array")
	if assert.Nil(t, err) {
		assert.ElementsMatch(t, []*Nested{
			{
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   242,
					},
				},
			},
			{
				isArray: true,
				array: []*Nested{
					{
						nested: map[string]*Nested{
							"value": {
								isValue: true,
								value:   "string in nested array",
							},
						},
					},
				},
			},
			{
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   1242,
					},
				},
			},
		}, array)
	}

	_, err = nested.GetArray("nested", "nested", "somekey")
	assert.EqualError(t, err, "nested.nested: key 'somekey' not found")

	_, err = nested.GetArray("nested", "nested")
	assert.EqualError(t, err, "nested.nested: is nested")

	nested = Nested{
		isArray: true,
		array: []*Nested{
			{
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   "string in nested array",
					},
				},
			},
		},
	}

	array, err = nested.GetArray()
	if assert.Nil(t, err) {
		assert.ElementsMatch(t, []*Nested{
			{
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   "string in nested array",
					},
				},
			},
		}, array)
	}

	nested = Nested{
		isValue: true,
		value:   42,
	}

	_, err = nested.GetArray()
	assert.EqualError(t, err, "is value")

	nested = Nested{}

	_, err = nested.GetArray()
	assert.EqualError(t, err, "is empty")
}

func Test_SetArray(t *testing.T) {
	nested := testNested()

	array := []*Nested{
		{
			isValue: true,
			value:   42,
		},
		{
			isValue: true,
			value:   "string",
		},
	}

	err := nested.SetArray(array, "array", "somekey")
	assert.EqualError(t, err, "array: is array")

	err = nested.SetArray(array, "nested", "nested", "value", "somekey")
	assert.EqualError(t, err, "nested.nested.value: is value")

	err = nested.SetArray(array, "somekey")
	if assert.Nil(t, err) {
		res, err := nested.GetArray("somekey")
		if assert.Nil(t, err) {
			assert.ElementsMatch(t, array, res)
		}
	}

	err = nested.SetArray(array, "nested", "nested", "somekey")
	if assert.Nil(t, err) {
		res, err := nested.GetArray("nested", "nested", "somekey")
		if assert.Nil(t, err) {
			assert.ElementsMatch(t, array, res)
		}
	}

	nested = Nested{}

	err = nested.SetArray(array)
	if assert.Nil(t, err) {
		res, err := nested.GetArray()
		if assert.Nil(t, err) {
			assert.ElementsMatch(t, array, res)
		}
	}
}

func Test_Delete(t *testing.T) {
	nested := testNested()

	err := nested.Delete()
	assert.EqualError(t, err, "keys list must contain at least one key")

	err = nested.Delete("somekey", "nested")
	assert.EqualError(t, err, "key 'somekey' not found")

	err = nested.Delete("nested", "array", "somekey")
	assert.EqualError(t, err, "nested.array: is array")

	err = nested.Delete("nested", "value", "somekey")
	assert.EqualError(t, err, "nested.value: is value")

	err = nested.Delete("array", "somekey")
	assert.EqualError(t, err, "array: is array")

	err = nested.Delete("value", "somekey")
	assert.EqualError(t, err, "value: is value")

	err = nested.Delete("nested", "nested", "value")
	if assert.Nil(t, err) {
		nestedNested, err := nested.Get("nested", "nested")
		if assert.Nil(t, err) {
			assert.True(t, nestedNested.IsEmpty())
		}
	}

	err = nested.Delete("nested", "value")
	if assert.Nil(t, err) {
		_, err := nested.Get("nested", "nested")
		assert.Nil(t, err)
		_, err = nested.Get("nested", "array")
		assert.Nil(t, err)
		_, err = nested.Get("nested", "value")
		assert.EqualError(t, err, "nested: key 'value' not found")
	}

	nested = Nested{
		isValue: true,
		value:   42,
	}
	err = nested.Delete("somekey")
	assert.EqualError(t, err, "is value")

	nested = Nested{
		isArray: true,
	}
	err = nested.Delete("somekey")
	assert.EqualError(t, err, "is array")

	nested = Nested{
		nested: map[string]*Nested{
			"value1": {
				isValue: true,
				value:   142,
			},
			"value2": {
				isValue: true,
				value:   242,
			},
		},
	}
	err = nested.Delete("somekey")
	assert.Nil(t, err)

	err = nested.Delete("value1")
	if assert.Nil(t, err) {
		assert.Equal(t,
			map[string]*Nested{
				"value2": {
					isValue: true,
					value:   242,
				},
			},
			nested.nested)
	}
}

func Test_ArrayAdd(t *testing.T) {
	nested := testNested()

	newNested := &Nested{
		isValue: true,
		value:   242,
	}

	err := nested.ArrayAdd(newNested, "array", "somekey")
	assert.EqualError(t, err, "array: is array")

	err = nested.ArrayAdd(newNested, "nested", "nested", "value")
	assert.EqualError(t, err, "nested.nested.value: is value")

	err = nested.ArrayAdd(newNested, "nested", "nested")
	assert.EqualError(t, err, "nested.nested: is nested")

	err = nested.ArrayAdd(newNested, "nested", "nested", "somekey")
	assert.EqualError(t, err, "nested.nested: key 'somekey' not found")

	err = nested.ArrayAdd(newNested, "array")
	if assert.Nil(t, err) {
		array, err := nested.GetArray("array")
		if assert.Nil(t, err) {
			assert.ElementsMatch(t, []*Nested{
				{
					isValue: true,
					value:   142,
				},
				{
					isValue: true,
					value:   "string in array",
				},
				{
					isValue: true,
					value:   242,
				},
			}, array)
		}
	}

	err = nested.ArrayAdd(newNested, "nested", "array")
	if assert.Nil(t, err) {
		array, err := nested.GetArray("nested", "array")
		if assert.Nil(t, err) {
			assert.ElementsMatch(t, []*Nested{
				{
					nested: map[string]*Nested{
						"value": {
							isValue: true,
							value:   242,
						},
					},
				},
				{
					isArray: true,
					array: []*Nested{
						{
							nested: map[string]*Nested{
								"value": {
									isValue: true,
									value:   "string in nested array",
								},
							},
						},
					},
				},
				{
					nested: map[string]*Nested{
						"value": {
							isValue: true,
							value:   1242,
						},
					},
				},
				{
					isValue: true,
					value:   242,
				},
			}, array)
		}
	}

	nested = Nested{
		isValue: true,
		value:   42,
	}

	err = nested.ArrayAdd(newNested)
	assert.Error(t, err, "is value")

	nested = Nested{}

	err = nested.ArrayAdd(newNested)
	assert.Error(t, err, "is nested")

	err = nested.ArrayAdd(newNested, "somekey")
	assert.Error(t, err, "is empty")

	nested = Nested{
		isArray: true,
		array: []*Nested{
			{
				isValue: true,
				value:   142,
			},
		},
	}

	err = nested.ArrayAdd(newNested, "somekey")
	assert.Error(t, err, "is array")

	err = nested.ArrayAdd(newNested)
	if assert.Nil(t, err) {
		array, err := nested.GetArray()
		if assert.Nil(t, err) {
			assert.ElementsMatch(t, []*Nested{
				{
					isValue: true,
					value:   142,
				},
				{
					isValue: true,
					value:   242,
				},
			},
				array)
		}
	}
}

func Test_ArrayAddValue(t *testing.T) {
	nested := testNested()

	err := nested.ArrayAddValue(242, "array")
	if assert.Nil(t, err) {
		array, err := nested.GetArray("array")
		if assert.Nil(t, err) {
			assert.ElementsMatch(t, []*Nested{
				{
					isValue: true,
					value:   142,
				},
				{
					isValue: true,
					value:   "string in array",
				},
				{
					isValue: true,
					value:   242,
				},
			}, array)
		}
	}
}

func Test_ArrayAddArray(t *testing.T) {
	nested := testNested()

	newArray := []*Nested{
		{
			isValue: true,
			value:   1,
		},
		{
			isValue: true,
			value:   2,
		},
		{
			isValue: true,
			value:   3,
		},
	}

	err := nested.ArrayAddArray(newArray, "array")
	if assert.Nil(t, err) {
		array, err := nested.GetArray("array")
		if assert.Nil(t, err) {
			assert.ElementsMatch(t, []*Nested{
				{
					isValue: true,
					value:   142,
				},
				{
					isValue: true,
					value:   "string in array",
				},
				{
					isArray: true,
					array: []*Nested{
						{
							isValue: true,
							value:   1,
						},
						{
							isValue: true,
							value:   2,
						},
						{
							isValue: true,
							value:   3,
						},
					},
				},
			}, array)
		}
	}
}

func Test_ArrayFindOne(t *testing.T) {
	nested := testNested()

	search := func(element *Nested) bool {
		return true
	}

	_, err := nested.ArrayFindOne(search, "somekey")
	assert.EqualError(t, err, "key 'somekey' not found")

	_, err = nested.ArrayFindOne(search, "nested", "somekey")
	assert.EqualError(t, err, "nested: key 'somekey' not found")

	_, err = nested.ArrayFindOne(search, "value")
	assert.EqualError(t, err, "value: is value")

	_, err = nested.ArrayFindOne(search)
	assert.EqualError(t, err, "is nested")

	_, err = nested.ArrayFindOne(search, "nested", "nested", "somekey")
	assert.EqualError(t, err, "nested.nested: key 'somekey' not found")

	_, err = nested.ArrayFindOne(search, "nested", "nested")
	assert.EqualError(t, err, "nested.nested: is nested")

	nested = Nested{
		isValue: true,
		value:   42,
	}

	_, err = nested.ArrayFindOne(search)
	assert.EqualError(t, err, "is value")

	nested = Nested{}

	_, err = nested.ArrayFindOne(search)
	assert.EqualError(t, err, "is empty")

	nested = testNested()

	search = func(element *Nested) bool {
		value, err := element.GetValue()
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v == 142
	}

	value, err := nested.ArrayFindOne(search, "array")
	if assert.Nil(t, err) {
		assert.Equal(t, 142, value.value.(int))
	}

	search = func(element *Nested) bool {
		value, err := element.GetValue()
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v > 0
	}

	value, err = nested.ArrayFindOne(search, "nested", "array")
	if assert.Nil(t, err) {
		assert.Nil(t, value)
	}

	search = func(element *Nested) bool {
		value, err := element.GetValue("value")
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v > 0
	}

	nestedValue, err := nested.ArrayFindOne(search, "nested", "array")
	if assert.Nil(t, err) && assert.NotNil(t, nestedValue) {
		value, _ := nestedValue.GetValue("value")
		assert.Equal(t, 242, value.(int))
	}

	search = func(element *Nested) bool {
		value, err := element.GetValue("value")
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v > 1000
	}

	nestedValue, err = nested.ArrayFindOne(search, "nested", "array")
	if assert.Nil(t, err) && assert.NotNil(t, nestedValue) {
		value, _ := nestedValue.GetValue("value")
		assert.Equal(t, 1242, value.(int))
	}
}

func Test_ArrayFindAll(t *testing.T) {
	nested := testNested()

	search := func(element *Nested) bool {
		return true
	}

	_, err := nested.ArrayFindAll(search, "somekey")
	assert.EqualError(t, err, "key 'somekey' not found")

	_, err = nested.ArrayFindAll(search, "nested", "somekey")
	assert.EqualError(t, err, "nested: key 'somekey' not found")

	_, err = nested.ArrayFindAll(search, "value")
	assert.EqualError(t, err, "value: is value")

	_, err = nested.ArrayFindAll(search)
	assert.EqualError(t, err, "is nested")

	_, err = nested.ArrayFindAll(search, "nested", "nested", "somekey")
	assert.EqualError(t, err, "nested.nested: key 'somekey' not found")

	_, err = nested.ArrayFindAll(search, "nested", "nested")
	assert.EqualError(t, err, "nested.nested: is nested")

	nested = Nested{
		isValue: true,
		value:   42,
	}

	_, err = nested.ArrayFindAll(search)
	assert.EqualError(t, err, "is value")

	nested = Nested{}

	_, err = nested.ArrayFindAll(search)
	assert.EqualError(t, err, "is empty")

	nested = testNested()

	search = func(element *Nested) bool {
		value, err := element.GetValue()
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v == 142
	}

	array, err := nested.ArrayFindAll(search, "array")
	if assert.Nil(t, err) && assert.Len(t, array, 1) {
		assert.Equal(t, 142, array[0].value.(int))
	}

	search = func(element *Nested) bool {
		return element.IsValue()
	}

	array, err = nested.ArrayFindAll(search, "nested", "array")
	if assert.Nil(t, err) {
		assert.Len(t, array, 0)
	}

	search = func(element *Nested) bool {
		value, err := element.GetValue("value")
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v > 0
	}

	array, err = nested.ArrayFindAll(search, "nested", "array")
	if assert.Nil(t, err) {
		assert.ElementsMatch(t, []*Nested{
			{
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   242,
					},
				},
			},
			{
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   1242,
					},
				},
			},
		},
			array)
	}

	search = func(element *Nested) bool {
		value, err := element.GetValue("value")
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v > 1000
	}

	array, err = nested.ArrayFindAll(search, "nested", "array")
	if assert.Nil(t, err) {
		assert.ElementsMatch(t, []*Nested{
			{
				nested: map[string]*Nested{
					"value": {
						isValue: true,
						value:   1242,
					},
				},
			},
		},
			array)
	}
}

func Test_ArrayDelete(t *testing.T) {
	nested := testNested()

	search := func(element *Nested) bool {
		return true
	}

	err := nested.ArrayDelete(search, "somekey")
	assert.EqualError(t, err, "key 'somekey' not found")

	err = nested.ArrayDelete(search, "nested", "somekey")
	assert.EqualError(t, err, "nested: key 'somekey' not found")

	err = nested.ArrayDelete(search, "value")
	assert.EqualError(t, err, "value: is value")

	err = nested.ArrayDelete(search)
	assert.EqualError(t, err, "is nested")

	err = nested.ArrayDelete(search, "nested", "nested", "somekey")
	assert.EqualError(t, err, "nested.nested: key 'somekey' not found")

	err = nested.ArrayDelete(search, "nested", "nested")
	assert.EqualError(t, err, "nested.nested: is nested")

	nested = Nested{
		isValue: true,
		value:   42,
	}

	err = nested.ArrayDelete(search)
	assert.EqualError(t, err, "is value")

	nested = Nested{}

	err = nested.ArrayDelete(search)
	assert.EqualError(t, err, "is empty")

	nested = testNested()

	search = func(element *Nested) bool {
		return element.IsValue()
	}

	err = nested.ArrayDelete(search, "array")
	if assert.Nil(t, err) {
		array, _ := nested.GetArray("array")
		assert.Len(t, array, 0)
	}

	search = func(element *Nested) bool {
		value, err := element.GetValue("value")
		if err != nil {
			return false
		}

		v, ok := value.(int)
		if !ok {
			return false
		}

		return v > 0
	}

	err = nested.ArrayDelete(search, "nested", "array")
	if assert.Nil(t, err) {
		array, _ := nested.GetArray("nested", "array")
		assert.ElementsMatch(t, []*Nested{
			{
				isArray: true,
				array: []*Nested{
					{
						nested: map[string]*Nested{
							"value": {
								isValue: true,
								value:   "string in nested array",
							},
						},
					},
				},
			},
		},
			array)
	}
}

func Test_FromObject(t *testing.T) {
	// детальное использование тестируется в [Test_FromJSONString].
	nested := FromObject(42)
	assert.True(t, nested.IsValue())

	nested = FromObject([]any{1, 2, 3})
	assert.True(t, nested.IsArray())

	nested = FromObject(map[string]any{"a": 1, "b": 2})
	assert.True(t, nested.IsNested())
	assert.False(t, nested.IsEmpty())
}

func Test_FromJSONString(t *testing.T) {
	nested := FromJSONString(`42`)
	assert.Equal(t,
		&Nested{
			isValue: true,
			value:   42,
		},
		nested,
	)

	nested = FromJSONString(`42.5`)
	assert.Equal(t,
		&Nested{
			isValue: true,
			value:   42.5,
		},
		nested,
	)

	nested = FromJSONString(`string`)
	assert.Equal(t,
		&Nested{
			isValue: true,
			value:   "string",
		},
		nested,
	)

	nested = FromJSONString(`true`)
	assert.Equal(t,
		&Nested{
			isValue: true,
			value:   true,
		},
		nested,
	)

	nested = FromJSONString(`[[4, 5]`)
	assert.Equal(t,
		&Nested{
			isValue: true,
			value:   `[[4, 5]`,
		},
		nested,
	)

	nested = FromJSONString(`{"string3": skjhgdf}`)
	assert.Equal(t,
		&Nested{
			isValue: true,
			value:   `{"string3": skjhgdf}`,
		},
		nested,
	)

	nested = FromJSONString(`{"str": "string", "number": 42}`)
	assert.Equal(t,
		&Nested{
			nested: map[string]*Nested{
				"number": {
					isValue: true,
					value:   42,
				},
				"str": {
					isValue: true,
					value:   "string",
				},
			},
		},
		nested,
	)

	nested = FromJSONString(`{"str": "string", "number": 1.5, "nested": {"str": "string2", "number": 2, "array": [7]}, "array": [3, "string3", [4, 5], {"number": 6}]}`)
	assert.Equal(t,
		&Nested{
			nested: map[string]*Nested{
				"number": {
					isValue: true,
					value:   1.5,
				},
				"str": {
					isValue: true,
					value:   "string",
				},
				"array": {
					isArray: true,
					array: []*Nested{
						{
							isValue: true,
							value:   3,
						},
						{
							isValue: true,
							value:   "string3",
						},
						{
							isArray: true,
							array: []*Nested{
								{
									isValue: true,
									value:   4,
								},
								{
									isValue: true,
									value:   5,
								},
							},
						},
						{
							nested: map[string]*Nested{
								"number": {
									isValue: true,
									value:   6,
								},
							},
						},
					},
				},
				"nested": {
					nested: map[string]*Nested{
						"str": {
							isValue: true,
							value:   "string2",
						},
						"number": {
							isValue: true,
							value:   2,
						},
						"array": {
							isArray: true,
							array: []*Nested{
								{
									isValue: true,
									value:   7,
								},
							},
						},
					},
				},
			},
		},
		nested,
	)

	nested = FromJSONString(`[3, "string3", [4, 5], {"number": 6}]`)
	assert.Equal(t,
		&Nested{
			isArray: true,
			array: []*Nested{
				{
					isValue: true,
					value:   3,
				},
				{
					isValue: true,
					value:   "string3",
				},
				{
					isArray: true,
					array: []*Nested{
						{
							isValue: true,
							value:   4,
						},
						{
							isValue: true,
							value:   5,
						},
					},
				},
				{
					nested: map[string]*Nested{
						"number": {
							isValue: true,
							value:   6,
						},
					},
				},
			},
		},
		nested,
	)
}

func Test_ToObject(t *testing.T) {
	// детальное использование тестируется в [Test_ToJSONString].
	nested := Nested{
		isValue: true,
		value:   42,
	}
	assert.Equal(t, 42, nested.ToObject().(int))

	nested = Nested{
		isArray: true,
		array: []*Nested{
			{
				isValue: true,
				value:   1,
			},
			{
				isValue: true,
				value:   2,
			},
			{
				isValue: true,
				value:   3,
			},
		},
	}

	array, ok := nested.ToObject().([]any)
	if assert.True(t, ok) {
		assert.Equal(t, 1, array[0].(int))
		assert.Equal(t, 2, array[1].(int))
		assert.Equal(t, 3, array[2].(int))
	}

	nested = Nested{
		nested: map[string]*Nested{
			"a": {
				isValue: true,
				value:   1,
			},
			"b": {
				isValue: true,
				value:   2,
			},
		},
	}

	kv, ok := nested.ToObject().(map[string]any)
	if assert.True(t, ok) {
		assert.Equal(t, 1, kv["a"].(int))
		assert.Equal(t, 2, kv["b"].(int))
	}
}

func Test_ToJSONString(t *testing.T) {
	nested := Nested{
		isValue: true,
		value:   42,
	}

	assert.Equal(t,
		"42",
		nested.ToJSONString(),
	)

	nested = Nested{
		isValue: true,
		value:   42.5,
	}

	assert.Equal(t,
		"42.5",
		nested.ToJSONString(),
	)

	nested = Nested{
		isValue: true,
		value:   true,
	}

	assert.Equal(t,
		"true",
		nested.ToJSONString(),
	)

	nested = Nested{
		isValue: true,
		value:   "string",
	}

	assert.Equal(t,
		"string",
		nested.ToJSONString(),
	)

	nested = Nested{}

	assert.Equal(t,
		"{}",
		nested.ToJSONString(),
	)

	nested = Nested{}
	nested.SetArray([]*Nested{})
	assert.Equal(t,
		"[]",
		nested.ToJSONString(),
	)

	nested = Nested{
		nested: map[string]*Nested{
			"number": {
				isValue: true,
				value:   42,
			},
			"str": {
				isValue: true,
				value:   "string",
			},
		},
	}

	nested = *FromJSONString(`{"str": "string", "number": 42}`)
	assert.Equal(t,
		`{"number":42,"str":"string"}`,
		nested.ToJSONString(),
	)

	nested = Nested{
		nested: map[string]*Nested{
			"number": {
				isValue: true,
				value:   1.5,
			},
			"str": {
				isValue: true,
				value:   "string",
			},
			"array": {
				isArray: true,
				array: []*Nested{
					{
						isValue: true,
						value:   3,
					},
					{
						isValue: true,
						value:   "string3",
					},
					{
						isArray: true,
						array: []*Nested{
							{
								isValue: true,
								value:   4,
							},
							{
								isValue: true,
								value:   5,
							},
						},
					},
					{
						nested: map[string]*Nested{
							"number": {
								isValue: true,
								value:   6,
							},
						},
					},
				},
			},
			"nested": {
				nested: map[string]*Nested{
					"str": {
						isValue: true,
						value:   "string2",
					},
					"number": {
						isValue: true,
						value:   2,
					},
					"array": {
						isArray: true,
						array: []*Nested{
							{
								isValue: true,
								value:   7,
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t,
		`{"array":[3,"string3",[4,5],{"number":6}],"nested":{"array":[7],"number":2,"str":"string2"},"number":1.5,"str":"string"}`,
		nested.ToJSONString(),
	)

	nested = Nested{
		isArray: true,
		array: []*Nested{
			{
				isValue: true,
				value:   3,
			},
			{
				isValue: true,
				value:   "string3",
			},
			{
				isArray: true,
				array: []*Nested{
					{
						isValue: true,
						value:   4,
					},
					{
						isValue: true,
						value:   5,
					},
				},
			},
			{
				nested: map[string]*Nested{
					"number": {
						isValue: true,
						value:   6,
					},
				},
			},
		},
	}

	assert.Equal(t,
		`[3,"string3",[4,5],{"number":6}]`,
		nested.ToJSONString(),
	)
}
