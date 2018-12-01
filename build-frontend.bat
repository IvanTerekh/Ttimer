cd frontend
START /WAIT npm run build
cd ..
xcopy frontend\build\index.html views\index.html /Y
xcopy frontend\build\static static /Y/E
