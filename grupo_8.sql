
-- TABELA PERFIL
CREATE TABLE perfil (
    perfil_id   SERIAL PRIMARY KEY,
    tipo_perfil VARCHAR(50)  NOT NULL,
    descricao   VARCHAR(255) NOT NULL
);

----
-- TABELA LOJAS
CREATE TABLE lojas (
    loja_id       SERIAL PRIMARY KEY,
    nome          VARCHAR(50)  NOT NULL,
    luc           VARCHAR(25)  NOT NULL UNIQUE,
    segmento      VARCHAR(50)  NOT NULL,
    data_registro TIMESTAMP    DEFAULT CURRENT_TIMESTAMP
);

----
-- TABELA USUÁRIOS
CREATE TABLE usuarios (
    usuario_id   SERIAL PRIMARY KEY,
    nome         VARCHAR(100) NOT NULL,
    departamento VARCHAR(50)  NOT NULL,
    cargo        VARCHAR(50)  NOT NULL,
    perfil_id INTEGER REFERENCES perfil(perfil_id)
);

----
-- TABELA OCORRÊNCIAS
CREATE TYPE area_rspnsvl AS ENUM ('Arq', 'Eng', 'Bri');
CREATE TYPE category_ocorrencias AS ENUM (
    'Manutenção',
    'Conservação',
    'Conservação de resíduos',
    'Arquitetura e Paisagismo',
    'Segurança',
    'Brigada',
    'Engenharia'
);

CREATE TABLE ocorrencias (
    ocorrencia_id    SERIAL PRIMARY KEY,
    loja_id          INT                  NOT NULL,
    area_responsavel area_rspnsvl         NOT NULL,
    categoria        category_ocorrencias NOT NULL,
    assunto          VARCHAR(255)         NOT NULL,
    descricao        VARCHAR(255)         NOT NULL,
    data_registro    TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loja_id) REFERENCES lojas(loja_id)
);

----
-- TABELA MULTAS
CREATE TYPE status_multa AS ENUM ('Pendente', 'Faturada', 'Cancelada');
CREATE TYPE category_multas AS ENUM (
    'Manutenção e Engenharia',
    'Arquitetura e Paisagismo',
    'Conservação e Resíduos',
    'Segurança e Brigada'
);

CREATE TABLE multas (
    multa_id      SERIAL PRIMARY KEY,
    ocorrencia_id INT             NOT NULL,
    categoria     category_multas NOT NULL,
    assunto       VARCHAR(255)    NOT NULL,
    valor_multa   NUMERIC(10, 2)  NOT NULL DEFAULT 0.00,
    status        status_multa    NOT NULL DEFAULT 'Pendente',
    FOREIGN KEY (ocorrencia_id) REFERENCES ocorrencias(ocorrencia_id)
);

----
-- TABELA VISTORIAS
CREATE TYPE category_vistoria AS ENUM ('Reformar loja', 'Novo Contrato');
CREATE TABLE vistorias (
    vistoria_id      SERIAL PRIMARY KEY,
    loja_id          INT               NOT NULL,
    usuario_id       INT               NOT NULL,
    area_responsavel area_rspnsvl      NOT NULL,
    categoria        category_vistoria NOT NULL,
    assunto          VARCHAR(255),
    descricao        VARCHAR(255)      NOT NULL,
    data_registro    TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loja_id)    REFERENCES lojas(loja_id),
    FOREIGN KEY (usuario_id) REFERENCES usuarios(usuario_id)
);

----
-- TABELA ANEXOS
CREATE TABLE anexos (
    anexo_id      SERIAL PRIMARY KEY,
    nome_arquivo  VARCHAR(255) NOT NULL,
    tipo_mime     VARCHAR(100) NOT NULL,
    tamanho_bytes INT          NOT NULL,
    conteudo      BYTEA        NOT NULL,
    ocorrencia_id INT REFERENCES ocorrencias(ocorrencia_id) ON DELETE CASCADE,
    multa_id      INT REFERENCES multas(multa_id)           ON DELETE CASCADE,
    vistoria_id   INT REFERENCES vistorias(vistoria_id)     ON DELETE CASCADE,
    CONSTRAINT check_vinculo_anx CHECK (
        (ocorrencia_id IS NOT NULL AND vistoria_id IS NULL     AND multa_id IS NULL    ) OR
        (ocorrencia_id IS NULL     AND vistoria_id IS NOT NULL AND multa_id IS NULL    ) OR
        (ocorrencia_id IS NULL     AND vistoria_id IS NULL     AND multa_id IS NOT NULL)
    )
);

CREATE TABLE auditoria (
    auditoria_id   SERIAL PRIMARY KEY,
    entidade_tipo  VARCHAR(50)  NOT NULL, -- Ex: 'ocorrencia', 'multa', 'vistoria'
    entidade_id    INT          NOT NULL,
    acao           VARCHAR(50)  NOT NULL, -- Ex: 'create', 'update', 'delete', 'status'
    usuario        VARCHAR(100) NOT NULL,
    campo          VARCHAR(50),           -- Opcional: qual campo foi editado
    valor_anterior TEXT,                  -- Opcional: valor antes da edição
    valor_novo     TEXT,                  -- Opcional: valor após a edição
    detalhes       TEXT,                  -- Opcional: mensagem de log
    data_hora      TIMESTAMP    DEFAULT CURRENT_TIMESTAMP
);