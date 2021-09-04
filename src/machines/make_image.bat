xcopy ..\contracts\machines .\contracts\machines\ /s /e
docker build --tag dum-machines-service .
rmdir contracts /s/q