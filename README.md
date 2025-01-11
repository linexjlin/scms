# A very simple cms for markdown/html files using pandoc

## build

```bash
go build -o scms.exe main.go
```

## run

```bat
scms.exe -hot "C:\Users\line\OneDrive\data\LearnEnglish" -hist "C:\Users\line\OneDrive\data\LearnEnglish\history" -title "学习英语" -listen "127.0.0.1:8080"
```

## 使用

- 访问 http://127.0.0.1:8080/
