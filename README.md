# NFF - Nota Fiscal Fácil

Este programa automatiza processos repetitivos que funcionários públicos que trabalham no SIAT fazem para gerar notas fiscais no site do Siare para produtores rurais.

## Público alvo

O NFF é destinado à funcionários públicos, que possuem o mínimo de experiência com excel, que trabalham no setor do SIAT. Produtores rurais sempre precisam emitir notas fiscais devido ao grande número de transferências de gado que fazem. E é a esses funcionários que geralmente recorrem. Este projeto tem o objetivo de facilitar a vida destes funcionários, agilizando o trabalho que fazem.

Como meu irmão trabalha nesse setor, ele deu a ideia, eu vi que era viável, e assim se deu.

## Detalhes Técnicos

As principais tecnologias utilizadas foram:

* `python`
* `selenium`
* `pyinstaller` (para gerar o executável)
* `pandas`
* `webdriver-manager` (para evitar ter que baixar manualmente o driver do navegador)
* `excel` (utilizado como base de dados)

## Como usar?

* O primeiro passo é preencher a base de dados, que é o arquivo `"db.xlsx"` (por favor, note que o nome desse arquivo não deve mudar).

* Para preencher ele corretamente, siga o exemplo mostrado no arquivo `"db.example.xlsx"`.

* Na aba de **"Nota Fiscal"** dentro do arquivo excel, cada linha representa uma nota fiscal que será emitida na próxima execução do programa. Sendo assim, certifique-se de verificar essa aba antes de iniciar.

* Na aba **"Dados de Produtos e Serviços NF"** dentro do arquivo excel, o campo **"NF"** serve para indicar a qual nota fiscal que aquele produto pertence. Por exemplo, se eu colocar **_"1"_**, isso significa que aquele produto **_se refere à primeira nota fiscal da aba "Nota Fiscal"_**.

## Como funciona?

* Quando o programa iniciar, ele vai ler a base de dados, fazer as verificações necessárias, abrir o site do _Siare_ e emitir todas as notas fiscais que estiverem no excel (`"db.xlsx"`).

* As notas fiscais baixadas serão salvas em uma pasta chamada `"docs"` dentro do mesmo diretório do programa.

* Caso o campo **"senha"** no excel estiver vazio na hora de fazer login na conta de algum produtor, será aberta uma janelinha onde você pode digitar a senha.

## IMPORTANTE

* Não altere a aba **"Dados das listas"** da base de dados (`"db.xlsx"`) sem ter consciência do que está fazendo. As informações contidas ali servem para popular os campos de seleção que existem nas outras abas.

* Sempre lembre de salvar suas alterações na base de dados antes de iniciar o programa.

* **Se você decidir usar o campo "senha" no excel, garanta que ninguém mais tenha acesso à ele**.

## Limitações Atuais

* Não é possível emitir notas fiscais em casos onde alguma entidade não possui inscrição estadual.

* Também não é possível escolher com qual nome o arquivo da nota fiscal será salvo.

## Planos para o futuro

- [ ] Lidar com casos de destinatário sem inscrição municipal
- [ ] Ao final da execução, mostrar as NFs feitas com sucesso e as que não foram
- [ ] Ter um modo de mudar o nome do arquivo da nota fiscal
