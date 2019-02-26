@echo off

go build "investigo.go"
mkdir "Investigo-win32"
move "Investigo.exe" "Investigo-win32"
copy "sites.json" "Investigo-win32\sites.json"
python build.py
