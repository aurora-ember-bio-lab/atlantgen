# 1) Установить Trivy, если ещё нет (пропусти, если уже стоит / есть в Amber-Shield)
winget install AquaSecurity.trivy

# 2) Просканировать aura-amber-saas
trivy fs --format json --output "$env:USERPROFILE\aura-amber-saas\Claude outputs\trivy-report.json" "$env:USERPROFILE\aura-amber-saas"

# 3) Показать краткий читаемый вывод тоже (для тебя)
trivy fs "$env:USERPROFILE\aura-amber-saas"
