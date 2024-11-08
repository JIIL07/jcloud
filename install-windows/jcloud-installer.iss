[Setup]
AppName=Jcloud Client
AppVersion=0.0.1
DefaultDirName={commonpf}\JcloudClient
DefaultGroupName=Jcloud Client
OutputDir=userdocs:Output
OutputBaseFilename=JcloudClientSetup
Compression=lzma
SolidCompression=yes
SetupIconFile=installer-icon.ico

[Files]
Source: "config/*.yaml"; DestDir: "{app}/config"; Flags: ignoreversion
Source: "internal/client/*.*"; DestDir: "{app}/src"; Flags: ignoreversion recursesubdirs createallsubdirs

[Run]
Filename: "{app}\bin\jcloud.exe"; Description: "Run Jcloud Client"; Flags: nowait postinstall skipifsilent

[Tasks]
Name: "addtopath"; Description: "Add Jcloud Client to PATH"; GroupDescription: "Additional icons:"; Flags: unchecked

[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "PATH"; ValueData: "{olddata};{app}\bin"; Tasks: addtopath

[Code]
function DownloadFileWithPowerShell(URL, Dest: string): Boolean;
var
  ResultCode: Integer;
  PowerShellCommand: string;
begin
  Result := False;
  
  // Формируем команду PowerShell с правильными кавычками
  PowerShellCommand := 'powershell -NoProfile -ExecutionPolicy Bypass -Command "Invoke-WebRequest -Uri ''' + URL + ''' -OutFile ''' + Dest + ''' -ErrorAction Stop"';

  // Запускаем PowerShell команду через cmd.exe
  if Exec('cmd.exe', '/C ' + PowerShellCommand, '', SW_HIDE, ewWaitUntilTerminated, ResultCode) then
  begin
    if ResultCode = 0 then
      Result := True
    else
      MsgBox('Error downloading file: ' + URL + '. PowerShell exit code: ' + IntToStr(ResultCode), mbError, MB_OK);
  end
  else
    MsgBox('Failed to start PowerShell for downloading: ' + URL, mbError, MB_OK);
end;

procedure CurStepChanged(CurStep: TSetupStep);
var
  BinaryPath: string;
begin
  if CurStep = ssInstall then
  begin
    BinaryPath := ExpandConstant('{app}\bin\');
    
    // Создание папки для бинарных файлов
    if not DirExists(BinaryPath) then
      ForceDirectories(BinaryPath);

    // Скачиваем interactive.exe
    if not DownloadFileWithPowerShell('https://jcloud.up.railway.app/static/binary/interactive.exe', BinaryPath + 'interactive.exe') then
      MsgBox('Error downloading interactive.exe', mbError, MB_OK);

    // Скачиваем jcloud.exe
    if not DownloadFileWithPowerShell('https://jcloud.up.railway.app/static/binary/jcloud.exe', BinaryPath + 'jcloud.exe') then
      MsgBox('Error downloading jcloud.exe', mbError, MB_OK);
  end;
end;

