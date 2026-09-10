; Aura Amber Installer - Inno Setup Script
; Compile with Inno Setup 6+

#define MyAppName "Aura Amber"
#define MyAppVersion "0.3.0"
#define MyAppPublisher "Aurora Ember Bio Lab Ltd."
#define MyAppURL "https://github.com/cargounetcom/aura-amber-saas"
#define MyAppExeName "aura-amber.exe"

[Setup]
AppId={{AURA-AMBER-0001-0001-000000000001}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
OutputDir=installer
OutputBaseFilename=aura-amber-setup-{#MyAppVersion}
Compression=lzma2/ultra64
SolidCompression=yes
WizardStyle=modern
WizardSizePercent=110
SetupIconFile=installer\icon.ico
UninstallDisplayIcon={app}\{#MyAppExeName}
LicenseFile=LICENSE
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"
Name: "quicklaunchicon"; Description: "{cm:CreateQuickLaunchIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; Main binary
Source: "aura-amber.exe"; DestDir: "{app}"; Flags: ignoreversion
; Docker files
Source: "docker-compose.yml"; DestDir: "{app}"; Flags: ignoreversion
Source: "Dockerfile"; DestDir: "{app}"; Flags: ignoreversion
; Worker
Source: "worker\*"; DestDir: "{app}\worker"; Flags: ignoreversion recursesubdirs createallsubdirs
; UI
Source: "ui\*"; DestDir: "{app}\ui"; Flags: ignoreversion recursesubdirs createallsubdirs
; API source (for building)
Source: "cmd\*"; DestDir: "{app}\cmd"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "internal\*"; DestDir: "{app}\internal"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "atlas\*"; DestDir: "{app}\atlas"; Flags: ignoreversion recursesubdirs createallsubdirs
; Config
Source: ".env.example"; DestDir: "{app}"; Flags: ignoreversion
; Docs
Source: "README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "LICENSE"; DestDir: "{app}"; Flags: ignoreversion
Source: "API.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "CONTRIBUTING.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "go.mod"; DestDir: "{app}"; Flags: ignoreversion
Source: "go.sum"; DestDir: "{app}"; Flags: ignoreversion

[Dirs]
Name: "{app}\data"; Permissions: users-modify

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Parameters: "open"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Parameters: "open"; Tasks: desktopicon
Name: "{userappdata}\Microsoft\Internet Explorer\Quick Launch\{#MyAppName}"; Filename: "{app}\{#AppExeName}"; Parameters: "open"; Tasks: quicklaunchicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Parameters: "status"; Description: "Check service status"; Flags: postinstall nowait
Filename: "{app}\{#MyAppExeName}"; Parameters: "open"; Description: "Open dashboard in browser"; Flags: postinstall nowait skipifsilent

[Code]
// Glass-style wizard modernization
procedure InitializeWizard;
begin
  WizardForm.BorderStyle := bsNone;
  WizardForm.Position := poScreenCenter;
end;
