package orm

import (
	"reflect"
	"slices"
	"strings"
)

func UpdateMap(existing, updated interface{}, skipFields []string) map[string]interface{} {
	// 初始化最终更新字段映射表
	updates := make(map[string]interface{})

	// 防御性检查：nil 值直接返回空 map
	if existing == nil || updated == nil {
		return updates
	}

	// 获取反射值，并支持指针或值类型
	existingValue := reflect.ValueOf(existing)
	updatedValue := reflect.ValueOf(updated)

	if existingValue.Kind() == reflect.Ptr {
		existingValue = existingValue.Elem()
	}
	if updatedValue.Kind() == reflect.Ptr {
		updatedValue = updatedValue.Elem()
	}

	// 类型检查：确保两个结构体类型一致
	if existingValue.Type() != updatedValue.Type() {
		return updates
	}

	// 获取原始结构体类型，用于遍历字段和读取标签
	existingType := existingValue.Type()
	// 遍历结构体所有字段
	for i := 0; i < existingValue.NumField(); i++ {
		// 获取当前遍历的结构体字段信息
		field := existingType.Field(i)
		fieldName := field.Name
		// 获取原始值和更新值的字段反射值
		existingFieldValue := existingValue.Field(i)
		updatedFieldValue := updatedValue.Field(i)
		// 深度对比两个字段值是否不同（支持切片、map、结构体等复杂类型）
		if !reflect.DeepEqual(existingFieldValue.Interface(), updatedFieldValue.Interface()) {
			// 获取 gorm 标签，用于解析数据库列名
			dbTag := field.Tag.Get("gorm")
			// 默认列名 = 结构体字段名（无 gorm:column 标签时使用）
			columnName := fieldName
			// 从 gorm 标签中解析 column:xxx 定义的数据库列名
			if dbTag != "" {
				// 按分号分割 gorm 标签（gorm 标签支持多配置用;分隔）
				for _, tag := range strings.Split(dbTag, ";") {
					if strings.HasPrefix(tag, "column:") {
						// 提取列名并替换默认值
						columnName = strings.TrimPrefix(tag, "column:")
						break
					}
				}
			}
			// 跳过指定不需要对比的字段
			if slices.Contains(skipFields, columnName) {
				continue
			}
			// 将变化的字段（数据库列名 => 新值）存入更新映射表
			updates[columnName] = updatedFieldValue.Interface()
		}
	}

	return updates
}
