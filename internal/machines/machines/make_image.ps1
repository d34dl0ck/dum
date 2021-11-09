Copy-Item ..\contracts\machines -Filter "*.*" -Recurse -Destination .\contracts\machines -Verbose -Container
docker build --tag dum-machines-service .
Remove-Item contracts -Recurse -Force -Verbose