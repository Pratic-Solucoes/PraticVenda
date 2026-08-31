-- Migração para instalações existentes: consolida autenticação no schema public.
BEGIN;

ALTER TABLE public.tb_usuarios_gestao
    ADD COLUMN IF NOT EXISTS username VARCHAR(100);

UPDATE public.tb_usuarios_gestao
SET username = 'usuario_' || id
WHERE username IS NULL OR username = '';

ALTER TABLE public.tb_usuarios_gestao
    ALTER COLUMN username SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_usuarios_gestao_username
    ON public.tb_usuarios_gestao (username);

ALTER TABLE public.tb_usuarios_gestao
    DROP COLUMN IF EXISTS id_empresa;

DROP TABLE IF EXISTS public.tb_usuarios_admin;
DROP TABLE IF EXISTS public.tb_empresas_gestao;

COMMIT;
