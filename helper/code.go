package helper

import (
	"PracticeProject/define"
	"io/ioutil"
	"os"
)

// CodeSave
// 保存代码
func CodeSave(code []byte) (string, string, error) {
	dirName := "code/code-user/" + GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", "", err
	}
	f, err := os.Create(path)
	f.Write(code)
	defer f.Close()
	return path, dirName, nil
}

// CodeDelete
// 删除代码
func CodeDelete(dir string) error {
	err := os.RemoveAll(dir)
	return err
}

// CheckGoCodeValid
// 检查golang代码的合法性(import包的限制)
func CheckGoCodeValid(path string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	code := string(b)
	for i := 0; i < len(code)-6; i++ {
		if code[i:i+6] == "import" {
			var flag byte
			for i = i + 7; i < len(code); i++ {
				if code[i] == ' ' {
					continue
				}
				flag = code[i]
				break
			}
			if flag == '(' {
				for i = i + 1; i < len(code); i++ {
					if code[i] == ')' {
						break
					}
					if code[i] == '"' {
						t := ""
						for i = i + 1; i < len(code); i++ {
							if code[i] == '"' {
								break
							}
							t += string(code[i])
						}
						if _, ok := define.ValidGolangPackageMap[t]; !ok {
							return false, nil
						}
					}
				}
			} else if flag == '"' {
				t := ""
				for i = i + 1; i < len(code); i++ {
					if code[i] == '"' {
						break
					}
					t += string(code[i])
				}
				if _, ok := define.ValidGolangPackageMap[t]; !ok {
					return false, nil
				}
			}
		}
	}
	return true, nil
}
