$ErrorActionPreference = "SilentlyContinue"
$Path = "$($ENV:windir)\System32\Tasks"
$null = New-PSDrive -PSProvider registry -root HKEY_CLASSES_ROOT -Name HKCR
Get-ChildItem -Path $Path -Recurse | Where-Object { ! $_.PSIsContainer } | ForEach-Object {
	$TaskName = $_.Name
	$TaskXML = [xml] (Get-Content $_.FullName)
	if($TaskXML.Task.Actions.ComHandler) {     
		$TaskTrigger = $TaskXML.Task.Triggers.OuterXML
		$TaskXML.Task.Actions.Exec.Command| ForEach-Object {

			$COM = $TaskXML.Task.Actions.ComHandler.ClassID
			$dll = (Get-ItemProperty -LiteralPath HKCR:\CLSID\$COM\InprocServer32).'(default)'
			$Out = New-Object PSObject
			$Out | Add-Member Noteproperty 'TaskName' $TaskName
			$Out | Add-Member Noteproperty 'CLSID' $COM
			$Out | Add-Member Noteproperty 'Dll' $dll
			$Out | Add-Member Noteproperty 'Logon' $False
			$null = $TaskXML.Task.InnerXml -match 'Context="(?<Context>InteractiveUsers|AllUsers|AnyUser)"'

			$IsUserContext = $False
			if ($Matches['Context']) { $IsUserContext = $True }
			$Out | Add-Member Noteproperty 'IsUserContext' $IsUserContext

			if($TaskTrigger.Contains('LogonTrigger')){
				$Out.Logon = $True
			}
			else{$Out.Logon = $False}
			
			$Context = $null
						  
			
			if ($Out.Logon) {
					$Out
			}
			
		}
	}
}
