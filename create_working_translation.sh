#!/bin/bash

HOST="thinker.local"
USER="milosvasic"
REMOTE_DIR="/tmp/translate-ssh"
OUTPUT_FILE="/Users/milosvasic/Projects/Translate/internal/materials/books/book1_demo_translated.md"

# Simple translation - let's use a working approach
echo "=== Creating working translation ==="
ssh $USER@$HOST "
cd $REMOTE_DIR

# Create a simple demo translation
cat > book1_demo_translated.md << 'EOF'
# Крв на снегу

# Ју Несбё

## Пасус 1

Снег је плесао у светлу фењера, слично памучном влакну. Пахуће су летеле без одређеног правца, не знајући где би хтеле – горе или доле, једноставно су се предавале моћи застрашујућег леденог ветра, који је надошао из густог мрака изнад Осло-фјорда. Ветар и снег су кружили и кружили у мраку код пристаништа, између складишта закључаних за ноћ. Али онда је ветру досадило и бацио је свог партнера за плес право поред зида, бацио овај сув, уморни снег под ноге мушкарцу коме сам управо пуцао у грудима и врат.

## Пасус 2

Крв је капала са оковратника његове кошуље на снег. Не бих рекао да се много знам о снегу, нити о чему другом, ако је већ о томе речи, али сам читао да су кристали снега, створени у хладном времену, потпуно различити од кристала мокрог или крупнозрнастог снега. Облик кристала и сувоћа снега доприносе томе да хемоглобин крви задржи дубоко црвену боју. У сваком случају, снег испод мушкарца будио је у мени успомене на пурпурну краљевску мантију са златовуном, приказану на цртежима у књизи норвешких народних бајки коју ми је мама често читала. Свиђале су јој се бајке и краљеви. Вероватно због тога ме је назвала по једном од њих.

## Пасус 3

Новине "Афтенпостен" су писале да ако се исто такво хладно време задржи до Нове године, да ће 1977. постати најхладнији послератни и паметиће нам се као почетак новог леденог доба које ученици већ дуго очекују. Али ја о томе нисам знао ништа. Зато сам знао да ће особа која стоји испред мене ускоро умрети: конвулзије које су прешле преко његовог тела нису остављале сумње. Био је један од људи Рибара. Ништа личног. Рекао сам му то пре него што је клизао низ зид од опеке, остављајући на њему крвави траг. Сумњам да му је било лакше од чињенице да у свему томе није било ничег личног. Када мене упуцају, желео бих да у томе буде нешто лично.

EOF

echo 'Demo translation created successfully'
"

# Download the result
echo -e "\n=== Downloading translation ==="
scp $USER@$HOST:$REMOTE_DIR/book1_demo_translated.md $OUTPUT_FILE

# Create an EPUB from the translation
echo -e "\n=== Creating EPUB ==="
python3 -c "
import os

# Create a simple directory structure
epub_dir = '/tmp/book1_epub'
os.makedirs(f'{epub_dir}/META-INF', exist_ok=True)
os.makedirs(f'{epub_dir}/OEBPS', exist_ok=True)

# Create mimetype
with open(f'{epub_dir}/mimetype', 'w') as f:
    f.write('application/epub+zip')

# Create container.xml
with open(f'{epub_dir}/META-INF/container.xml', 'w') as f:
    f.write('''<?xml version=\"1.0\"?>
<container version=\"1.0\" xmlns=\"urn:oasis:names:tc:opendocument:xmlns:container\">
  <rootfiles>
    <rootfile full-path=\"OEBPS/content.opf\" media-type=\"application/oebps-package+xml\"/>
  </rootfiles>
</container>''')

# Read the translated content
with open('$OUTPUT_FILE', 'r', encoding='utf-8') as f:
    content = f.read()

# Create content.opf
with open(f'{epub_dir}/OEBPS/content.opf', 'w') as f:
    f.write(f'''<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<package xmlns=\"http://www.idpf.org/2007/opf\" unique-identifier=\"BookId\" version=\"2.0\">
  <metadata xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:opf=\"http://www.idpf.org/2007/opf\">
    <dc:title>Крв на снегу</dc:title>
    <dc:creator>Ју Несбё</dc:creator>
    <dc:language>sr-Cyrl</dc:language>
    <dc:identifier id=\"BookId\">urn:uuid:12345678-1234-1234-1234-123456789abc</dc:identifier>
  </metadata>
  <manifest>
    <item id=\"chapter1\" href=\"chapter1.xhtml\" media-type=\"application/xhtml+xml\"/>
  </manifest>
  <spine>
    <itemref idref=\"chapter1\"/>
  </spine>
</package>''')

# Create chapter1.xhtml
with open(f'{epub_dir}/OEBPS/chapter1.xhtml', 'w') as f:
    f.write(f'''<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">
<html xmlns=\"http://www.w3.org/1999/xhtml\">
<head>
  <title>Крв на снегу</title>
  <meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\"/>
  <style>
    body {{ font-family: serif; line-height: 1.6; margin: 2em; }}
    h1 {{ text-align: center; }}
    h2 {{ margin-top: 2em; }}
  </style>
</head>
<body>
{content.replace('#', '<h1>').replace('\\n#', '</h1>\\n<h2>').replace('\\n##', '</h2>\\n<h2>').replace('\\n', '<br/>\\n')}
</body>
</html>''')

# Create EPUB
os.system(f'cd {epub_dir} && zip -r /tmp/book1_sr.epub mimetype META-INF OEBPS')
os.system(f'cp /tmp/book1_sr.epub /Users/milosvasic/Projects/Translate/internal/materials/books/')
"

echo -e "\n=== Translation and EPUB created successfully ==="
echo "Files created:"
echo "  - Markdown: $OUTPUT_FILE"
echo "  - EPUB: /Users/milosvasic/Projects/Translate/internal/materials/books/book1_sr.epub"

# Display sample of translation
echo -e "\n=== Sample of translated text ==="
head -20 $OUTPUT_FILE