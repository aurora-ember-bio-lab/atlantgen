# 1) Удалить случайную вложенную пустую папку (внутри только .git, ничего не потеряется)
rm -rf "/c/Users/ellag/aura-amber-saas/aura-amber-saas"

# 2) Перейти в настоящий корень проекта
cd "/c/Users/ellag/aura-amber-saas"

# 3) Проверить, что git видит правильные файлы (ui/next-app и .gitignore)
git status

# 4) Закоммитить дашборд
git add .
git commit -m "Add Next.js migration dashboard UI"

# 5) Если ещё не подключён remote — подключить и запушить
# git remote add origin git@github.com:cargounetcom/aura-amber-saas.git
git push -u origin main
