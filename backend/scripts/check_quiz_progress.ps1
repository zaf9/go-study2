#!/usr/bin/env pwsh
# Quiz completion progress tracker
# Usage: ./check_quiz_progress.ps1

$ErrorActionPreference = "Stop"

$quizDataRoot = Join-Path -Path $PSScriptRoot -ChildPath "..\quiz_data"

if (-not (Test-Path $quizDataRoot)) {
    Write-Host "ERROR: Quiz data directory not found: $quizDataRoot" -ForegroundColor Red
    exit 1
}

# Define chapter manifest
$manifest = @{
    "lexical_elements" = @("comments", "keywords", "identifiers", "tokens", "integers", "floats", "strings", "semicolons", "imaginary", "runes")
    "constants" = @("boolean", "rune", "integer", "floating_point", "complex", "string")
    "variables" = @("variable_declarations", "short_declarations", "blank_identifier", "type_conversions", "zero_value")
    "types" = @("slice", "map", "interface", "struct", "pointer", "channel", "function", "array", "string", "numeric", "boolean", "type_definitions", "type_aliases", "type_parameters", "type_constraints", "type_inference", "type_unification", "underlying_types", "type_identity", "method_sets", "type_assertions")
}

$totalChapters = 0
$completedChapters = 0
$totalQuestions = 0
$targetMin = 1200
$targetMax = 2000

$topicStats = @{}

Write-Host ""
Write-Host "==== Quiz Bank Progress Report ====" -ForegroundColor Cyan
Write-Host "Generated: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor Gray
Write-Host ""

foreach ($topic in $manifest.Keys | Sort-Object) {
    $chapters = $manifest[$topic]
    $topicCompleted = 0
    $topicQuestions = 0
    
    foreach ($chapter in $chapters) {
        $totalChapters++
        $filePath = Join-Path -Path (Join-Path -Path $quizDataRoot -ChildPath $topic) -ChildPath "$chapter.yaml"
        
        if (Test-Path $filePath) {
            $content = Get-Content $filePath -Raw
            if ($content -match 'questions:\s*\n((?:\s+-\s+id:.*\n(?:\s+.*\n)*)+)') {
                $questionCount = ([regex]::Matches($content, '^\s*-\s+id:', [System.Text.RegularExpressions.RegexOptions]::Multiline)).Count
                
                if ($questionCount -ge 30) {
                    $completedChapters++
                    $topicCompleted++
                    $topicQuestions += $questionCount
                    $totalQuestions += $questionCount
                }
            }
        }
    }
    
    $topicStats[$topic] = @{
        Total = $chapters.Count
        Completed = $topicCompleted
        Questions = $topicQuestions
    }
}

# Overall statistics
$completionRate = if ($totalChapters -gt 0) { [math]::Round(($completedChapters / $totalChapters) * 100, 1) } else { 0 }
$progressBar = "#" * [math]::Floor($completionRate / 5)
$emptyBar = "-" * (20 - $progressBar.Length)

Write-Host "Total chapters: $totalChapters" -ForegroundColor White
Write-Host "Completed: $completedChapters ($completionRate%)" -ForegroundColor Green
Write-Host "Pending: $($totalChapters - $completedChapters) ($([math]::Round(100 - $completionRate, 1))%)" -ForegroundColor Yellow
Write-Host "Total questions: $totalQuestions / Target $targetMin-$targetMax" -ForegroundColor White
Write-Host "Progress: [$progressBar$emptyBar] $completionRate%" -ForegroundColor Cyan
Write-Host ""

# By topic
Write-Host "--- By Topic ---" -ForegroundColor Cyan
$topicNames = @{
    "lexical_elements" = "Lexical Elements"
    "constants" = "Constants"
    "variables" = "Variables"
    "types" = "Types"
}

foreach ($topic in @("lexical_elements", "constants", "variables", "types")) {
    $stat = $topicStats[$topic]
    $name = $topicNames[$topic]
    $rate = if ($stat.Total -gt 0) { [math]::Round(($stat.Completed / $stat.Total) * 100, 1) } else { 0 }
    
    $status = if ($rate -eq 100) { "[OK]" } elseif ($rate -ge 50) { "[PARTIAL]" } else { "[TODO]" }
    
    $msg = "[$name] {0}/{1} completed ({2}%) {3} ({4} questions)" -f $stat.Completed, $stat.Total, $rate, $status, $stat.Questions
    Write-Host $msg -ForegroundColor $(
        if ($rate -eq 100) { "Green" } elseif ($rate -ge 50) { "Yellow" } else { "Red" }
    )
}

Write-Host ""

# Quality targets
Write-Host "--- Quality Targets ---" -ForegroundColor Cyan

if ($totalQuestions -ge $targetMin -and $totalQuestions -le $targetMax) {
    Write-Host "[OK] Question count on target ($totalQuestions in range $targetMin-$targetMax)" -ForegroundColor Green
} elseif ($totalQuestions -lt $targetMin) {
    $deficit = $targetMin - $totalQuestions
    Write-Host "[WARN] Insufficient questions, need $deficit more to reach minimum target $targetMin" -ForegroundColor Yellow
} else {
    Write-Host "[OK] Sufficient questions ($totalQuestions exceeds target)" -ForegroundColor Green
}

if ($completionRate -eq 100) {
    Write-Host "[OK] All chapters completed!" -ForegroundColor Green
} else {
    $pending = $totalChapters - $completedChapters
    Write-Host "[WARN] Still have $pending chapters pending" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "==== End of Report ====" -ForegroundColor Cyan
Write-Host ""

# Return status code (0 = all done, 1 = pending)
if ($completionRate -eq 100 -and $totalQuestions -ge $targetMin) {
    exit 0
} else {
    exit 1
}


# 返回状态码 (0 = 全部完成, 1 = 有待完成)
if ($completionRate -eq 100 -and $totalQuestions -ge $targetMin) {
    exit 0
} else {
    exit 1
}
