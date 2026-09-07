# Диагностика: что реально отслеживает git-репозиторий в домашней папке

echo "--- top-level ---"
git -C ~ rev-parse --show-toplevel

echo "--- сколько файлов отслеживается ---"
git -C ~ ls-files | wc -l

echo "--- первые 20 отслеживаемых файлов ---"
git -C ~ ls-files | head -20

echo "--- есть ли .gitignore и что в нём ---"
cat ~/.gitignore 2>/dev/null | head -30

echo "--- отслеживается ли node_modules (не должен) ---"
git -C ~ ls-files -- node_modules | head -5

echo "--- статус (без untracked, только то что git реально видит как изменения) ---"
git -C ~ status --porcelain=v1 -uno
