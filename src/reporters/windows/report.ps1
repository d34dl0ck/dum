$session = New-Object -ComObject Microsoft.Update.Session
$searcher = $session.CreateupdateSearcher()
# This will be used after successful testing of ps1 output
# $updates = @($searcher.Search("IsHidden=0 and IsInstalled=0 and CategoryIDs contains 'E0789628-CE08-4437-BE74-2495B842F43B' or IsHidden=0 and IsInstalled=0 and CategoryIDs contains 'E6CF1350-C01B-414D-A61F-263D14D133B4' or IsHidden=0 and IsInstalled=0 and CategoryIDs contains '0FA1201D-4330-4FA8-8AE9-B877473B6441' or IsHidden=0 and IsInstalled=0 and CategoryIDs contains 'CD5FFD1E-E932-4E3A-BF74-18BF0B1BBD83'").Updates)
$updates = @($searcher.Search("IsHidden=0 and IsInstalled=0").Updates)
$dtoSet = [System.Collections.Generic.List[object]]::new()

foreach ($update in $updates) {
    $rawSeverity = $update.MsrcSeverity
    $severity = 0

    if ($rawSeverity -eq "Low") {
        $severity = 1
    } elseif ($rawSeverity -eq "Important"){
        $severity = 2
    } elseif ($rawSeverity -eq "Critical") {
        severity = 3
    }

    $utcNow = [DateTime]::UtcNow
    $durationOfMissing = $utcNow - ([DateTime]$update.LastDeploymentChangeTime)
    $duration = [math]::Round($durationOfMissing.TotalMinutes).ToString() + "m"

    $dto = [pscustomobject]@{
        UpdateId = $update.Identity.UpdateId;
        Severity = $severity;
        Duration = $duration;
    }

    $dtoSet.Add($dto)
}

$machineName = $env:computername
$reportUri = "http://localhost:3000/api/v1/machines/" + $machineName + "/report"

Invoke-WebRequest -Uri $reportUri -Method POST -Body ($dtoSet|ConvertTo-Json) -ContentType "application/json" -UseBasicParsing