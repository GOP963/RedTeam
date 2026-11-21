$ByteArray = [System.IO.File]::ReadAllBytes("C:\Users\NEO\Desktop\RavinAcademy\RedTeam-01\10\HTTPMalware\Persist\x64\Release\DLL.dll");
$Base64String = [System.Convert]::ToBase64String($ByteArray);
$Base64String