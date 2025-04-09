## Cenário
O Magalu tem o desafio de desenvolver uma plataforma de comunicação. Você foi
escolhido(a) para iniciar o desenvolvimento da primeira sprint.

## Objetivo
O objetivo nesta sprint é o de prover 3 serviços (endpoints) relativos ao envio de comunicação
da empresa.

1. Agendamento de envio de comunicação;
2. Consulta do envio da comunicação;
3. Cancelamento do envio da comunicação.

## Requisitos
Criar um endpoint que receba uma solicitação de agendamento de envio de comunicação (1);

Este endpoint precisa ter, no mínimo, os seguintes campos:
- Data/Hora para o envio
- Destinatário
- Mensagem a ser entregue

Os possíveis formatos de comunicação que podem ser enviados são:
- email, sms, push e whatsapp

Neste momento, precisamos deste canal de entrada para realizar o agendamento do envio, ou
seja, este endpoint (1).

O envio em si não será desenvolvido nesta etapa: você não precisa se preocupar com a
implementação do envio propriamente dito.

Para esta sprint, ficou decidido que a solicitação do agendamento do envio da comunicação
será salva no banco de dados. Portanto, assim que receber a solicitação do agendamento do
envio (1), ela deverá ser salva no banco de dados.

Pense com atenção nessa estrutura do banco. Apesar de não ser você quem vai realizar o
envio, essa estrutura já precisa estar pronta para que um colega seu não precise alterar nada
quando for desenvolver esta funcionalidade.

A preocupação no momento do envio será a de enviar e alterar o status do registro no banco
de dados. Contemple este ponto na sua modelagem da base de dados.

Deve ter um endpoint para consultar o status do agendamento de envio de comunicação (2). O
agendamento será feito no endpoint (1) e a consulta será feita por este outro endpoint (2).

Deve ter um endpoint para cancelar um agendamento de envio de comunicação (3).

## Observações e Orientações Gerais
Temos preferência por desenvolvimento na linguagem Java, Python ou Node, mas pode ser
usada qualquer outra linguagem. Apenas nos explique o porquê da sua escolha.

Utilize um dos bancos de dados abaixo:
- MySQL
- PostgreSQL

As APIs deverão seguir o modelo RESTFul com formato padrão JSON.

Faça testes unitários, foque em criar uma suite de testes bem organizada e automatizada.

Siga o que considera como boas práticas de programação.

A criação da base de dados e das suas tabelas fica a seu critério de como será feito, seja via
script, aplicação, etc.

Seu desafio deve ser enviado como repositório GIT público (Github, Gitlab, Bitbucket), com
commits pequenos e bem descritos. O seu repositório deve estar com um modelo de licença
de código aberto. Não envie nenhum arquivo além do próprio código e sua documentação.

Tome cuidado para não enviar imagens, vídeos, áudio, binários etc.
Siga boas práticas de desenvolvimento, de qualidade e de gestão e versionamento de código.

Oriente os avaliadores sobre como instalar, testar e executar seu código: pode ser um
README dentro do projeto.
Opcionalmente, forneça instruções ou um script e as configurações para rodar o projeto
localmente via Docker (Docker Composer). Este ponto será um bônus na avaliação.

Iremos avaliar seu desafio de acordo com a posição e o nível que você está se candidatando.

Agradecemos muito sua disposição de participar do nosso processo seletivo e desejamos que
você se divirta e que tenha boa sorte! :o)
