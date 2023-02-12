<img src="./assets/icon.ico" width="10%" align="left" alt="NFF main icon">

# NFF - Nota Fiscal Fácil

Este programa automatiza processos repetitivos que funcionários públicos (ou qualquer pessoa na verdade) fazem para gerar notas fiscais no site do [Siare](https://www2.fazenda.mg.gov.br/sol/) para produtores rurais.

## Download

Clique [aqui](https://drive.google.com/file/d/1so-2FLdHQxLCb8YIMhBXDHJMtxAYycwF/view?usp=share_link) para baixar o programa. Ele vem compactado no formato `.rar`. Basta extrair os arquivos normalmente. Não existe nenhum processo de instalação. As instruções de como usar estão logo abaixo.

## Sobre

Produtores rurais sempre precisam emitir notas fiscais devido ao grande número de transferências de gado e outros produtos que fazem. Geralmente eles recorrem à funcionários públicos para isso. Este projeto tem o objetivo de facilitar a vida destes funcionários, agilizando o seu trabalho.

Mas serve perfeitamente para qualquer pessoa que tenha o mínimo de afinidade com excel, afinal, esses funcionários públicos apenas realizam login na conta dos prórpios produtores para emitir notas fiscais.

Como meu irmão trabalha nesse setor, ele deu a ideia, eu vi que era viável, e assim se deu.

## Detalhes Técnicos

As principais tecnologias utilizadas foram:

* `python`
* `selenium`
* `pandas`
* `excel` (utilizado como base de dados)

## Como usar?

* O primeiro passo é preencher a base de dados, que é o arquivo `"db.xlsx"` (por favor, note que o nome desse arquivo não deve mudar).

* Para preencher ele corretamente, siga o exemplo mostrado no arquivo `"db.example.xlsx"`.

* Na aba de **"Nota Fiscal"** dentro do arquivo excel, cada linha representa uma nota fiscal que será emitida na próxima execução do programa. Sendo assim, certifique-se de verificar essa aba antes de iniciar.

* Ainda dentro da aba **"Nota Fiscal"**, os campos `"remetente"` e `"destinatário"` devem ser preenchidos com o `"cpf/cnpj"` das entidades correspondentes.

* Na aba **"Dados de Produtos e Serviços NF"** dentro do arquivo excel, o campo `"NF"` serve para indicar a qual nota fiscal que aquele produto pertence. Por exemplo, se eu colocar `1`, isso significa que aquele produto se refere à **_primeira nota fiscal da aba "Nota Fiscal"_**.

* Para iniciar, dê dois cliques no arquivo `NFF.exe`.

## Como funciona?

* Quando o programa iniciar, ele vai ler a base de dados, fazer as verificações necessárias, abrir o site do _Siare_, fazer login na conta do `"remetente"` e emitir todas as notas fiscais que estiverem no excel.

* As notas fiscais baixadas serão salvas em uma pasta chamada `"docs"` dentro do mesmo diretório do programa.

* Caso o campo `"senha"` no excel estiver vazio na hora de fazer login na conta de algum produtor, será aberta uma janelinha onde você pode digitar a senha.

## IMPORTANTE

* Não altere a aba **"Dados das listas"** da base de dados sem ter consciência do que está fazendo. As informações contidas ali servem para popular os campos de seleção que existem nas outras abas.

* Sempre lembre de salvar suas alterações na base de dados antes de iniciar o programa.

* **Se você decidir usar o campo `"senha"` no excel, garanta que ninguém mais tenha acesso à ele**.

## Limitações Atuais

* Não é possível escolher com qual nome o arquivo da nota fiscal será salvo.

## Próximos passos

- [x] Lidar com casos de destinatário sem inscrição estadual
- [ ] Ao final da execução, mostrar as NFs feitas com sucesso e as que não foram
- [ ] Ter um modo de mudar o nome do arquivo da nota fiscal
- [ ] Poder referenciar entidades na coluna `"remetente"` e `"destinatário"` tanto por `"cpf/cnpj"` como por `"inscrição estadual"`

## Considerações finais

Acredito que tenha ficado um projeto bem massinha. Obviamente existem pontos a ser melhorados, e também é muito provável que haja uma forma melhor de se automatizar esse processo, caso seja possível se comunicar com a API do governo, ao invés de ter que depender do front deles. Mas no momento isso tudo já é incrível, e a ideia é continuar evoluindo.
