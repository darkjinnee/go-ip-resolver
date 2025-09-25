import re
from urllib.parse import urlparse
import socket
import argparse

# Парсер аргументов
parser = argparse.ArgumentParser(description="Извлечение уникальных доменов из HTML и резолвинг их в IP")
parser.add_argument('-f', '--file', required=True, help="Путь к HTML файлу")
args = parser.parse_args()

file_path = args.file
unique_domains = set()

# Регулярка для поиска href и src
link_pattern = re.compile(r'(?:href|src)=["\'](.*?)["\']', re.IGNORECASE)

# Чтение файла
with open(file_path, 'r', encoding='utf-8') as f:
    content = f.read()

links = link_pattern.findall(content)

# Извлечение доменов
for link in links:
    parsed = urlparse(link)
    if parsed.netloc:
        unique_domains.add(parsed.netloc)

# Резолвинг доменов в IP
for domain in sorted(unique_domains):
    try:
        ip = socket.gethostbyname(domain)
        print(f"{domain} => {ip}")
    except socket.gaierror:
        print(f"{domain} => не удалось резолвить")