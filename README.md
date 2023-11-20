<img src="./web/internal/static/icons/icon.ico" width="10%" align="left" alt="NFF main icon">

# NFF - Nota Fiscal Fácil

Este projeto tem o objetivo de automatizar processos repetitivos que funcionários públicos (ou qualquer pessoa na verdade) fazem para requerer/baixar/cancelar notas fiscais estaduais em um site do governo chamado [Siare](https://www2.fazenda.mg.gov.br/sol/) para produtores rurais.

## Demonstração

### [Clique aqui para ver o projeto em produção](https://nff-web-dev.up.railway.app/)

> A nota fiscal emitida no vídeo é apenas para fins demonstrativos

[Vídeo da web]()

## Sobre

Produtores rurais sempre precisam emitir notas fiscais devido ao grande número de transferências de gado e outros produtos que fazem. As vezes também precisam cancelar notas. E todo final de ano precisam calcular o balanço de entrada e saída em notas fiscais de venda. Geralmente eles recorrem à funcionários públicos e contadores para isso. Este projeto tem o objetivo de facilitar a vida destes profissionais, agilizando o seu trabalho.

Mas serve perfeitamente para qualquer pessoa, afinal esses funcionários públicos apenas realizam login na conta dos próprios produtores para realizar essas tarefas.

Como meu irmão trabalha nesse setor, ele deu a ideia, eu vi que era viável, e assim se deu.

## Próximos passos

- [x] Lidar com casos de destinatário sem inscrição estadual
- [x] Ao final da execução, mostrar as operações feitas com sucesso e as que não foram
- [x] Possibilitar **cancelamento** de notas fiscais
- [x] Possibilitar **impressão** de notas fiscais isoladamente
- [x] Calcular **métricas** de entrada e saída em determinado período
- [x] Histórico de operações
- [x] Sistema de usuários (já implementado, porém sujeito a otimizações)
- [ ] Paginação de histórico de requerimentos
- [ ] Filtragem/Pesquisa no histórico de requerimentos
- [ ] Possibilitar o uso de modelos de nota fiscal, preenchendo automaticamente os campos (parcialmente implementado)
- [ ] Cancelamento de NF através do protocolo
- [ ] Login/Cadastro com Google
- [ ] Bom uso de cache
