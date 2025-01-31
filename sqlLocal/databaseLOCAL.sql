
ANTES DE RODAR O SCRIPT VC TEM QUE TROCAR OS XXX PELO NUMERO
DA EMPRESA QUE VC ESTÁ CRIANDO
TABELA EMPRESA
TABELA Usuário

CREATE DATABASE IF NOT EXISTS databaseLOCAL;
USE databaseLOCAL;


CREATE TABLE EMPRESA (
    ID_EMPRESA INT  AUTO_INCREMENT PRIMARY KEY, 
    CDEMP           VARCHAR(5)   NOT NULL, 
    NOME            VARCHAR(50)  DEFAULT "",
    DATA_SISTEMA    DATE,
    FANTASIA        VARCHAR(50)  DEFAULT "",
    ENDERECO        VARCHAR(50)  DEFAULT "",
    CIDADE          VARCHAR(50)  DEFAULT "",
    PAIS            VARCHAR(2)   DEFAULT "",
    COD_REG_EMP     VARCHAR(50)  DEFAULT "",
    FONE            VARCHAR(50)  DEFAULT "",
    CELULAR         VARCHAR(50)  DEFAULT "",
    PARAMETROS      VARCHAR(200) DEFAULT "",
    PY_AUTORIZADO   VARCHAR(200) DEFAULT "",   
    PY_PRE_FIXO     VARCHAR(200) DEFAULT "",
    PY_NR_BOL       INT DEFAULT 0,
    GRUPO           int DEFAULT 0,



    CRIADOEM   TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Campo criadoEm adicionado
) ENGINE=INNODB;



INSERT INTO EMPRESA (
    CDEMP,
    DATA_SISTEMA
) VALUES (
    "00XXX",
    CURDATE()
);





CREATE TABLE USUARIOS(
    ID int auto_increment primary key,
    ID_HOSTWEB  int DEFAULT 0,
    AGUARDANDO_SYNC VARCHAR(1) DEFAULT "",
    nome varchar(50) not null,
    NICK varchar(20) not null,
    EMAIL varchar(50) not null unique,
    SENHA varchar(100) not null,
    NIVEL int not null,
    CDEMP varchar(5) not null,
    
    CRIADOEM timestamp default current_timestamp()
) ENGINE=INNODB;

CADASTRAR UM USUARIO COM A ROTAS

troca numero CDEMP e apagar essa linha na hora de executar

insert into usuarios (nome, nick, email, senha,CDEMP, NIVEL)
values
("Usuário 00xxx", 
"usuario00xxx", 
"usuario00xxx@gmail.com",
 "$2a$10$0iGYlKCAYTyJV/vC6nLGgeWFwD6AhSkWLsVRO/.M4lNK8OtIkfggy",
 "00xxx", 
 100);



CREATE TABLE NIVEL_ACESSO(
    ID int auto_increment primary key,
    CODIGO varchar(50) not null unique,
    NOME_ACESSO varchar(50) not null unique,
    nivel int not null,
    CRIADOEM timestamp default current_timestamp()
) ENGINE=INNODB;


CREATE TABLE `CLIENTE` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `ID_HOSTWEB` int DEFAULT 0,
  `AGUARDANDO_SYNC` VARCHAR(1) DEFAULT "",
  `CNPJ_CPF` varchar(18) DEFAULT "",
  `NOME` varchar(50) DEFAULT "",
  `FANTASIA` varchar(50) DEFAULT "",
  `ENDERE` varchar(50) DEFAULT "",
  `CIDADE` varchar(50) DEFAULT "",
  `BAIRRO` varchar(30) DEFAULT "",
  `CEP` varchar(10) DEFAULT "",
  `UF` varchar(2) DEFAULT "",
  `TELE1` varchar(20) DEFAULT "",
  `TELE2` varchar(20) DEFAULT "",
  `TELE3` varchar(20) DEFAULT "",
  `CEL1` varchar(20) DEFAULT "",
  `CEL2` varchar(20) DEFAULT "",
  `FAX` varchar(20) DEFAULT "",
  `CONTATO` varchar(50) DEFAULT "",
  `CONTA_TEL1` varchar(20) DEFAULT "",
  `CONTA_TEL2` varchar(20) DEFAULT "",
  `EMAIL` varchar(50) DEFAULT "",
  `INSC_EST_RG` varchar(15) DEFAULT "",
  `TAXA` float DEFAULT 0,
  `PLACA` varchar(10) DEFAULT "",
  `NUMERO` varchar(5) DEFAULT "",
  `CIDIBGE` varchar(8) DEFAULT "",
  `COMPLE` varchar(100) DEFAULT "",
  `ID_CLIENTE_IFOOD` varchar(128) DEFAULT "",
  `DT_NASC` date DEFAULT NULL,
  `STATUS` varchar(1) DEFAULT 'A',
  `CRIADOEM` timestamp default current_timestamp(),
  PRIMARY KEY (`ID`),
  KEY `CLIENTE_IDX1` (`CNPJ_CPF`),
  KEY `CLIENTE_IDX2` (`NOME`),
  KEY `CLIENTE_IDX3` (`TELE1`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE `CLIENTE` AUTO_INCREMENT = 11125;


CREATE TABLE `PRODUTO` (
    `ID`             int NOT NULL AUTO_INCREMENT,
    `ID_HOSTWEB`     int DEFAULT 0,
    `ID_HOSTLOCAL`   int DEFAULT 0,
    `AGUARDANDO_SYNC` VARCHAR(1) DEFAULT "",
    `CODM`           varchar(20) DEFAULT '' UNIQUE,
    `DES1`           varchar(60) DEFAULT '',
    `DES2`           varchar(60) DEFAULT '',
    `PV`             FLOAT DEFAULT 0,
    `PV2`            FLOAT DEFAULT 0,
    `IMPRE`          varchar(2) DEFAULT '',
    `FORNE`          DOUBLE PRECISION DEFAULT 0,
    `CST1`           FLOAT DEFAULT 0,
    `CST2`           FLOAT DEFAULT 0,
    `STOCK1`         FLOAT DEFAULT 0,
    `STOCK2`         FLOAT DEFAULT 0,
    `STOCK3`         FLOAT DEFAULT 0,
    `STOCK4`         FLOAT DEFAULT 0,
    `COMPOSI`        varchar(1) DEFAULT '',
    `MINIMO`         FLOAT DEFAULT 0,
    `PS_L`           FLOAT DEFAULT 0,
    `PS_BR`          FLOAT DEFAULT 0,
    `COM_1`          FLOAT DEFAULT 0,
    `COM_2`          FLOAT DEFAULT 0,
    `BARRA`          varchar(20) DEFAULT '',
    `OB1`            varchar(1) DEFAULT '',
    `OB2`            varchar(1) DEFAULT '',
    `OB3`            varchar(1) DEFAULT '',
    `OB4`            varchar(1) DEFAULT '',
    `OB5`            varchar(1) DEFAULT '',
    `OBS1`           varchar(1) DEFAULT '',
    `OBS2`           varchar(1) DEFAULT '',
    `OBS3`           varchar(1) DEFAULT '',
    `OBS4`           varchar(1) DEFAULT '',
    `OBS5`           varchar(15) DEFAULT '',
    `STATUS`         varchar(1) DEFAULT '',
    `UND`            varchar(5) DEFAULT '',
    `CARDAPIO`       varchar(1) DEFAULT '',
    `GRUPO`          varchar(2) DEFAULT '',
    `FISCAL`         varchar(1) DEFAULT '',
    `FISCAL1`        varchar(1) DEFAULT '',
    `BAIXAD`         FLOAT DEFAULT 0,
    `SALDO_AN`       FLOAT DEFAULT 0,
    `CUSTO_CM`       FLOAT DEFAULT 0,
    `CUSTO_MD_AN`    FLOAT DEFAULT 0,
    `CUSTO_MD_AT`    FLOAT DEFAULT 0,
    `CNAE`           varchar(10) DEFAULT '',
    `CD_SERV`        varchar(10) DEFAULT '',
    `COMISSAO`       FLOAT DEFAULT 0,
    `CODPISCOFINS`   varchar(10) DEFAULT '',
    `SERIAL`         varchar(10) DEFAULT '',
    `QPCX`           INTEGER DEFAULT 0,
    `PV3`            FLOAT DEFAULT 0,
    `NCM`            varchar(15) DEFAULT '',
    `EAN`            varchar(15) DEFAULT '',
    `CSON`           varchar(3) DEFAULT '',
    `ORIGEN`         varchar(1) DEFAULT '',
    `ALICOTA`        FLOAT DEFAULT 0,
    `SERVICO`        varchar(1) DEFAULT '',
    `DES2_COMPLE`    varchar(300) DEFAULT '',
    `ST2`            varchar(2) DEFAULT '',
    `ST_PROMO1`      varchar(2) DEFAULT '',
    `PC_PROMO`       FLOAT DEFAULT 0,
    `PROM1_HI`       TIME DEFAULT '00:00:00',
    `PROM1_HF`       TIME DEFAULT '00:00:00',
    `PROM1_SEMANA`   varchar(2) DEFAULT '',
    `PROMO_CONSUMO`  FLOAT DEFAULT 0,
    `PROMO_PAGAR`    FLOAT DEFAULT 0,
    `CRIADOEM`       timestamp DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (`ID`),
    KEY `PRODUTO_IDX1` (`CODM`),
    KEY `PRODUTO_IDX2` (`DES1`),
    KEY `PRODUTO_IDX3` (`DES2`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



ALTER TABLE `PRODUTO` AUTO_INCREMENT = 10125;






CREATE TABLE CX_RECEB (
    ID           INT  AUTO_INCREMENT PRIMARY KEY,  -- Chave primária com auto incremento
    ID_HOSTWEB   int DEFAULT 0,
    STATUS       VARCHAR(1) DEFAULT "L", 
    DATA         DATE,
    ID_CLI       INT   DEFAULT 0, 
    ID_USER      INT   DEFAULT 0,   
    TOTAL        FLOAT DEFAULT 0,
    TROCO        FLOAT DEFAULT 0,
    MESA         INT DEFAULT 0,
    NR_PESSOAS   INT DEFAULT 0,
    CRIADOEM   TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Campo criadoEm adicionado
) ENGINE=INNODB;

CREATE INDEX CX_RECEB_IDX1 ON CX_RECEB (DATA);
CREATE INDEX CX_RECEB_IDX2 ON CX_RECEB (ID_CLI);
CREATE INDEX CX_RECEB_IDX3 ON CX_RECEB (ID_USER);

CREATE TABLE CX_RECEB_TIPO (
    ID           INT  AUTO_INCREMENT PRIMARY KEY,  -- Chave primária com auto incremento
    ID_CX_RECEB  int DEFAULT 0,
    ID_HOSTWEB   int DEFAULT 0,
    ID_TIPO_REC  INT DEFAULT 0,
    MOEDA_NAC    FLOAT DEFAULT 0,
    MOEDA_EXT    FLOAT DEFAULT 0



) ENGINE=INNODB;


CREATE TABLE TIPO_REC(
    ID           INT  AUTO_INCREMENT PRIMARY KEY,  -- Chave primária com auto incremento
    NOME         VARCHAR(30) DEFAULT "",
    CAMBIO       FLOAT DEFAULT 0,
    FT_CONV      VARCHAR(1),
    STATUS       VARCHAR(1) DEFAULT "A"

) ENGINE=INNODB;



     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Real BR ',   1, '*' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Dolar US ',   5.8, '*' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Guarani PY ',  1750, '/' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Peso ARG ',   140, '/' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Euro ',   5.8, '*' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Cartão de Credito ',   0, '' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Cartão de Debito ',   0, '' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Vale Refeição ',   0, '' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Pix ',   0, '' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'QR COD PY ',   0, '' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Desconto',   0, '' );  
     INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV) VALUES (   'Bonificação ',   0, '' );  







CREATE TABLE `IMPRESSORAS` (
    `ID`              int NOT NULL AUTO_INCREMENT PRIMARY KEY,  
    `ID_HOSTLOCAL`    int DEFAULT 0,    
    `AGUARDANDO_SYNC` VARCHAR(1) DEFAULT "",
    `COD_IMP`         VARCHAR(2)  DEFAULT "",
    `NOME`            VARCHAR(20) DEFAULT "",
    `END_IMP`         VARCHAR(50)  DEFAULT "",
    `END_SER`         VARCHAR(50)  DEFAULT ""
 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


insert into IMPRESSORAS (COD_IMP, NOME) values
("A", "COZINHA"),
("B", "COZINHA 2"),
("C", "COZINHA 3"),
("D", "BEBIDAS"),
("E", "FECHAMENTO"),
("F", "DELIVERY"),
("G", "VENDA DIRETA"),
("H", "OUTRAS 1"),
("I", "OUTRAS 2"),
("J", "OUTRAS 3"),
("K", "OUTRAS 4"),
("L", "OUTRAS 5"),
("M", "OUTRAS 6"),
("N", "OUTRAS 7"),
("O", "OUTRAS 8"),
("P", "OUTRAS 9"),
("Q", "OUTRAS 10"),
("R", "OUTRAS 11"),
("S", "OUTRAS 12"),
("T", "OUTRAS 13"),
("U", "OUTRAS 14"),
("V", "OUTRAS 15"),
("X", "OUTRAS 16"),
("Z", "OUTRAS 17");


CREATE TABLE `GRUPOS` (
    `ID`             int NOT NULL AUTO_INCREMENT PRIMARY KEY,  
    `ID_HOSTWEB`     int DEFAULT 0,    
    `ID_HOSTLOCAL`     int DEFAULT 0,   
    `AGUARDANDO_SYNC` VARCHAR(1) DEFAULT "",   
    `COD_GP`         VARCHAR(2)  DEFAULT "",
    `NOME`           VARCHAR(20) DEFAULT "",
    `TIPO`           VARCHAR(1)  DEFAULT "",
    `CONTA_CONTABIL` VARCHAR(6)  DEFAULT ""
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE VENDA (
    `ID`              int NOT NULL AUTO_INCREMENT PRIMARY KEY,  
    `CODM`            varchar(20) DEFAULT "",
    `STATUS`          varchar(1) DEFAULT "A",
    `IMPRESSORA`      varchar(2) DEFAULT "",  
    `MONITOR_PRD`     varchar(1) DEFAULT "",  
    `STATUS_TP_VENDA` varchar(1) DEFAULT "", 
    `STATUS_TP_HARD`  varchar(1) DEFAULT "", 

    `MESA`            int DEFAULT 0,
    `CELULAR`         VARCHAR(20) DEFAULT "",
    `CPF_CLIENTE`     VARCHAR(30) DEFAULT "",
    `NOME_CLIENTE`    VARCHAR(30) DEFAULT "",
    `ID_CLIENTE`      int DEFAULT 0,
    `PV`              float DEFAULT 0,
    `DESCONTO`        float DEFAULT 0,
    `PV_PROM`         float DEFAULT 0,
    `QTD`             float DEFAULT 0,
    `ID_USER`         INT DEFAULT 0,
    `NICK`            varchar(20) DEFAULT "",
    `DATA`            DATE,

    `OBS`             VARCHAR(100) DEFAULT "",
    `OBS2`            VARCHAR(100) DEFAULT "",
    `OBS3`            VARCHAR(100) DEFAULT "",
    `STATUS_PGTO`     VARCHAR(1) DEFAULT "",
    `DATA_IFOOD`      DATE DEFAULT NULL,
    `STATUS_IFOOD`    varchar(10) DEFAULT "",
    `ID_IFOOD`        VARCHAR(50) DEFAULT "",
    `COMPLEM_CODM`    VARCHAR(50) DEFAULT "", 
    `CHAVE`           VARCHAR(100) DEFAULT "",  
    `CRIADOEM`        timestamp default current_timestamp()

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE VENEXT (
    `ID`              int NOT NULL AUTO_INCREMENT PRIMARY KEY,  
    `ID_VENDA`        int DEFAULT 0,    
    `STATUS`          varchar(1) DEFAULT "",
    `CODM`            varchar(20) DEFAULT "",
    `ID_HOSTWEB`      int DEFAULT 0,  
    `ID_CX_RECEB`     int DEFAULT 0,
    `STATUS_TP_VENDA` varchar(1) DEFAULT "", 
    `STATUS_TP_HARD`  varchar(1) DEFAULT "", 
    
    `MESA`            int DEFAULT 0,
    `CELULAR`         VARCHAR(20) DEFAULT "",
    `CPF_CLIENTE`     VARCHAR(30) DEFAULT "",
    `NOME_CLIENTE`    VARCHAR(30) DEFAULT "",
    `ID_CLIENTE`      int DEFAULT 0,
   
    `PV`              float DEFAULT 0,
    `PV_PROM`         float DEFAULT 0,
    `QTD`             float DEFAULT 0,
    `ID_USER`         INT DEFAULT 0,
    `NICK`            varchar(20) DEFAULT "",
    `DATA`            DATE,
   

    `STOCK`           varchar(1) DEFAULT "",  
    `CUSTO`           float DEFAULT 0,
 
 
    `CRIADOEM`        timestamp default current_timestamp()

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE MESA (
    `ID`               int NOT NULL AUTO_INCREMENT PRIMARY KEY,  
    `MESA_CARTAO`      int not null unique,
    `STATUS`           varchar(1) DEFAULT "L",
    `ID_USER`          INT DEFAULT 0,
    `ID_CLI`           INT DEFAULT 0,
    `NICK`             varchar(50) DEFAULT "", 
 
    `ABERTURA`         timestamp,
    `QTD_PESSOAS`       INT DEFAULT 0,
    `TURISTA`          VARCHAR(1) DEFAULT "",
    `CELULAR`          VARCHAR(20) DEFAULT "",
    `APELIDO`          VARCHAR(20) DEFAULT ""

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;




CREATE TABLE `PARAMETROS` (
    `ID`                int NOT NULL AUTO_INCREMENT,
    `NOME`              VARCHAR(100),
    `STATUS`            VARCHAR(1),
    `LIMITE`            FLOAT DEFAULT 0,
     PRIMARY KEY (`ID`)
    

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;






     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Adm -Moeda Principal (R)eal (G)uarani (D)olar (P)eso AR  ',   'R'  );  
     INSERT INTO PARAMETROS (NOME, STATUS,LIMITE) VALUES (   'Adm - Taxa de serviço',   'N',10  );  
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Comanda produção - Imprimir Mix - todas as impressoras',   'N'  );  
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Comanda produção - Imprimir Cancelamento Impressora de Origem',   'N'  );   
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Caixa - Permitir receber como CORTESIA',   'N'  );  
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Caixa - Permitir receber como CONTAS A RECEBER',   'N'  );  
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Adm - Tipo de recebimento (M)esas (C)artões',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Caixa - Receber Mesa/Cartão sem imprimir ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Conta - Imprimir Comprovante de Pgto ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Comanda produção - Imprimir Lançamento direto Local',   'N'  );  
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Adm - No recebimento, Trazer valor total em sua Moeda Principal ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Conta - Conta Resumida (produtos igual, soma as quantidades) ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Comanda produção - Venda Direta não imprimir ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Comanda produção - Imprimir registro de vendas*',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Adm - Ocultar moedas Extrangeira na Conta',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Conta - Imprimir Taxa de Serviços separada ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Adm -Mostrar opção Numero de pessoas na mesa ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Adm - Mostrar opção TURISTA ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Conta - Venda Direta - Imprimir 2 vias ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Conta - Delivery -Imprimir 2 vias  ',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Conta - Mostrar Forma de Pagamentos no Delivery',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Conta - Mostrar nome da empresa na CONTA',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Comanda produção -01- IMPRIMIR*',   'S'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Comanda produção -02- SEPARAR LANÇAMENTO*',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Comanda produção -03- Imprimir em fichas individuais*',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS,LIMITE) VALUES (   'Adm - Limite de Cartão',   'N',1000  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Enviar Pedido - PESQUISA CLIENTE*',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Enviar Pedido - CPF*',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Enviar Pedido - NOME/APELIDO*',   'N'  ); 
     INSERT INTO PARAMETROS (NOME, STATUS) VALUES (   'Enviar Pedido - CELULAR*',   'N'  ); 
 

