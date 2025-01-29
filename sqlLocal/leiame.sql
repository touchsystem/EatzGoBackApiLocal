


CADASTRAR UM USUARIO COM A ROTAS

--primeiro insere um usaurio manualmente ex via postman

insert into usuarios (nome, nick, email, senha,CDEMP, NIVEL)
values
("Usuário 00199", "usuario00199", "usuario00199@gmail.com", "$2a$10$0iGYlKCAYTyJV/vC6nLGgeWFwD6AhSkWLsVRO/.M4lNK8OtIkfggy","00199", 0)


na rota \login-sync PUT aplica o json abaixo
{
  "email": "usuario00199@gmail.com",
  "nome": "usuario00199",
  "nick": "usuario00199",
  "cdEmp": "00199"

}


trocar id da empresa

usar ROTAS

localhost:5001/login-sync 

vai retorar no TOKEN-SYNC 

