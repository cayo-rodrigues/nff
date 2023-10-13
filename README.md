<img src="./web/internal/static/icons/icon.ico" width="10%" align="left" alt="NFF main icon">

# NFF - Nota Fiscal Fácil

Este projeto tem o objetivo de automatizar processos repetitivos que funcionários públicos (ou qualquer pessoa na verdade) fazem para requerer/baixar/cancelar notas fiscais estaduais em um site do governo chamado [Siare](https://www2.fazenda.mg.gov.br/sol/) para produtores rurais.

## Demonstração

> A nota fiscal emitida no vídeo é apenas para fins demonstrativos

https://user-images.githubusercontent.com/87717182/222856710-af0801b3-294c-43a8-8b9b-b5f0760a2e14.mp4

## Sobre

Produtores rurais sempre precisam emitir notas fiscais devido ao grande número de transferências de gado e outros produtos que fazem. As vezes também precisam cancelar notas. E todo final de ano precisam calcular o balanço de entrada e saída em notas fiscais de venda. Geralmente eles recorrem à funcionários públicos para isso. Este projeto tem o objetivo de facilitar a vida destes funcionários, agilizando o seu trabalho.

Mas serve perfeitamente para qualquer pessoa, afinal esses funcionários públicos apenas realizam login na conta dos próprios produtores para realizar essas tarefas.

Como meu irmão trabalha nesse setor, ele deu a ideia, eu vi que era viável, e assim se deu.

## Próximos passos

- [x] Lidar com casos de destinatário sem inscrição estadual
- [x] Ao final da execução, mostrar as NFs feitas com sucesso e as que não foram
- [x] Ter um modo de mudar o nome do arquivo da nota fiscal
- [x] Poder referenciar entidades na coluna `"remetente"` e `"destinatário"` tanto por `"cpf/cnpj"` como por `"inscrição estadual"`
- [x] Possibilitar **cancelamento** de notas fiscais
- [x] Aba de Histórico de notas fiscais emitidas/canceladas
- [ ] Possibilitar o uso de modelos de nota fiscal, preenchendo automaticamente os campos (já implementado, porém pode melhorar)

### Projeto no insomnia

Clique [aqui](https://drive.google.com/file/d/1wk0HeMX07f_M2HsvvOuSvlhfHrX85BKa/view?usp=sharing) para baixar.
