/*
测试数据需要使用MySQL, 并在MySQL中导入以下数据
-- ----------------------------
-- Table structure for student
-- ----------------------------
DROP TABLE IF EXISTS `student`;
CREATE TABLE `student`  (
  `stuid` int(11) NOT NULL,
  `stuname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `stubirth` datetime(0) NOT NULL,
  `stutel` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `stuaddr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `stuphoto` longblob NULL,
  PRIMARY KEY (`stuid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of student
-- ----------------------------
INSERT INTO `student` VALUES (1002, '郭靖', '1980-02-02 00:00:00', NULL, NULL, NULL);
INSERT INTO `student` VALUES (1003, '黄蓉', '1982-03-03 00:00:00', NULL, '成都市二环路南四段123号', NULL);
*/
package db_file

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"os"
	"testing"
)

func TestDB_Init(t *testing.T) {
	type fields struct {
		args *Arguments
		conn *sql.DB
	}
	type args struct {
		args *Arguments
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// 正确
		{name: "current",
			fields: fields{
				args: &Arguments{
					Url: "my://test:test@123@10.12.1.236:3307/test",
				}},
			args: args{
				args: &Arguments{
					Url: "my://test:test@123@10.12.1.236:3307/test",
				},
			},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				args: tt.fields.args,
				conn: tt.fields.conn,
			}
			// 初始化
			if err := d.Init(d.args); (err != nil) != tt.wantErr {
				t.Errorf("DB.NewDupKeyHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// 测试数据库连接
func TestDB_Connect(t *testing.T) {
	type fields struct {
		args *Arguments
		conn *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// 正确测试案例
		{
			name: "current",
			fields: fields{
				args: &Arguments{
					Url: "my://test:test!123@10.12.1.236:3307/test",
				},
			},
			wantErr: false,
		},
		// 错误测试案例
		{
			name: "error",
			fields: fields{
				args: &Arguments{
					Url: "my://test:test!123@10.12.1.236:3308/test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				args: tt.fields.args,
				conn: tt.fields.conn,
			}
			if err := d.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("DB.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// 测试执行查询
func TestDB(t *testing.T) {
	type fields struct {
		args *Arguments
		conn *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// 正确案例
		{name: "current",
			fields: fields{
				args: &Arguments{
					Url:                 "my://test:test!123@10.12.1.236:3307/test",
					ExecSql:             "select * from student",
					SaveFile:            "./result.txt",
					TagAll:              false,
					Tag:                 `'`,
					FromEncoding:        "UTF8",
					ToEncoding:          "UTF8",
					TagExcludeFieldType: "INT",
					EmptyVal:            `\N`,
					Separator:           ";",
					Newline:             "\n",
				},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				args: tt.fields.args,
				conn: tt.fields.conn,
			}
			// 初始化
			if err := d.Init(d.args); (err != nil) != tt.wantErr {
				t.Errorf("DB.NewDupKeyHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := d.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("DB.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := d.WriteFile(); (err != nil) != tt.wantErr {
				t.Errorf("DB.WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 测试列名称结果
			sColumnsName := []string{"stuid", "stuname", "stubirth", "stutel", "stuaddr", "stuphoto"}
			rColumnsName := d.GetColumnsName()
			if sColumnsName[0] != rColumnsName[0] {
				t.Errorf("test GetColumnsName error, got: %s, execpt: %s", sColumnsName[0], rColumnsName[0])
			}
			if sColumnsName[1] != rColumnsName[1] {
				t.Errorf("test GetColumnsName error, got: %s, execpt: %s", sColumnsName[1], rColumnsName[1])
			}
			if sColumnsName[2] != rColumnsName[2] {
				t.Errorf("test GetColumnsName error, got: %s, execpt: %s", sColumnsName[2], rColumnsName[2])
			}
			if sColumnsName[3] != rColumnsName[3] {
				t.Errorf("test GetColumnsName error, got: %s, execpt: %s", sColumnsName[3], rColumnsName[3])
			}
			if sColumnsName[4] != rColumnsName[4] {
				t.Errorf("test GetColumnsName error, got: %s, execpt: %s", sColumnsName[4], rColumnsName[4])
			}
			if sColumnsName[5] != rColumnsName[5] {
				t.Errorf("test GetColumnsName error, got: %s, execpt: %s", sColumnsName[5], rColumnsName[5])
			}
			// 测试获取数据类型
			sColumnsType := []string{"INT", "VARCHAR", "DATETIME", "CHAR", "VARCHAR", "BLOB"}
			rColumnsType := d.GetColumnsType()
			if sColumnsType[0] != rColumnsType[0] {
				t.Errorf("test GetColumnsType error, got: %s, execpt: %s", sColumnsType[0], rColumnsType[0])
			}
			if sColumnsType[1] != rColumnsType[1] {
				t.Errorf("test GetColumnsType error, got: %s, execpt: %s", sColumnsType[1], rColumnsType[1])
			}
			if sColumnsType[2] != rColumnsType[2] {
				t.Errorf("test GetColumnsType error, got: %s, execpt: %s", sColumnsType[2], rColumnsType[2])
			}
			if sColumnsType[3] != rColumnsType[3] {
				t.Errorf("test GetColumnsType error, got: %s, execpt: %s", sColumnsType[3], rColumnsType[3])
			}
			if sColumnsType[4] != rColumnsType[4] {
				t.Errorf("test GetColumnsType error, got: %s, execpt: %s", sColumnsType[4], rColumnsType[4])
			}
			if sColumnsType[5] != rColumnsType[5] {
				t.Errorf("test GetColumnsType error, got: %s, execpt: %s", sColumnsType[5], rColumnsType[5])
			}
			// 验证whetherAddTag
			if d.whetherAddTag("int") != false {
				t.Error("test whetherAddTag error, got: false, expect: true")
			}

			// 验证关闭连接
			if err := d.Close(); (err != nil) != tt.wantErr {
				t.Errorf("DB.Close(() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 验证写入文件结果是否正确
			var buf bytes.Buffer
			_, _ = buf.WriteString("1002;'郭靖';'1980-02-02 00:00:00';'\\N';'\\N';'\\N'\n1003;'黄蓉';'1982-03-03 00:00:00';'\\N';'成都市二环路南四段123号';'\\N'\n")
			sData := buf.String()
			rByte, err := ioutil.ReadFile("./result.txt")
			if err != nil {
				t.Error("read result file error, error message: ", err)
			}
			if sData != string(rByte[:]) {
				t.Errorf("test result file error, got: %q, expect: %q", sData, rByte[:])
			}
			_ = os.Remove("./result.txt")
		})
	}
}
