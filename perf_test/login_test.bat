@echo off
Remove-Item login_targets.json -Force
go run .\generate_users.go

echo Running Vegeta login test...
if not exist results mkdir results

Get-Content ./login_targets.json | vegeta attack -format=json -rate=1000 -duration=60s -connections=100 -max-connections=1000 -max-workers=1000 -output=results/login_result_1k_RPS.bin
vegeta report -type=text results/login_result_1k_RPS > results/login_report_1k_RPS.txt
vegeta plot results/login_result_1k_RPS > results/login_result_plot_1k_RPS.html
.\results\login_result_plot_1k_RPS.html
