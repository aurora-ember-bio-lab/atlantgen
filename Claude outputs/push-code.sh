# Выполнять в Git Bash

# 1) удалить пустую вложенную папку-дубль
rm -rf "/c/Users/ellag/aura-amber-saas/aura-amber-saas"

# 2) перейти в корень проекта
cd "/c/Users/ellag/aura-amber-saas"

# 3) закоммитить всё (UI на Next.js 16 + .gitignore)
git add .
git commit -m "Add Next.js migration dashboard UI"

# 4) запушить на GitHub
git push -u origin main
