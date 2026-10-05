		CONVENTIONAL COMMITS (intentad hacerlo en **`inglés`**):

**feat**: Añades una nueva funcionalidad al usuario. Ejemplo:  `feat: añadir login con Google`

**fix**: Corriges un error o bug. Ejemplo: `fix: corregir cálculo de IVA en el carrito`

**build**: Cambios que afectan al sistema de compilación o dependencias externas (npm, gradle, webpack...). Ejemplo: `build: actualizar versión de webpack a 5.x`

**chore**: Tareas de mantenimiento que no afectan a la lógica ni al usuario final (configuración, scripts, limpieza). Ejemplo: `chore: actualizar .gitignore`

**refactor**: Cambias el código sin alterar su comportamiento externo (ni añades funcionalidad ni arreglas bugs). Ejemplo: `refactor: simplificar función de validación`

**style**: Cambios de formato que no afectan a la lógica (espacios, comas, indentación, punto y coma). Ejemplo: `style: aplicar formato con prettier`

**docs**: Cambios exclusivamente en documentación. Ejemplo: `docs: actualizar README con instrucciones de instalación`

**test**: Añades o corriges tests. No modifica código de producción. Ejemplo: `test: añadir tests unitarios para el servicio de pagos`

**perf**: Cambios de código enfocados a mejorar el rendimiento. Ejemplo: `perf: optimizar consulta SQL de usuarios`

**ci**: Cambios en la configuración de integración continua (GitHub Actions, Jenkins, GitLab CI). Ejemplo: `ci: añadir workflow de despliegue automático`

**revert**: Deshaces un commit anterior. Ejemplo: `revert: revertir feat: login con Google`



	BRANCH NAMING (intentad hacerlo en **`inglés`**):

**feature/** (o feat/): Para desarrollar una nueva funcionalidad. Ejemplo: `feature/add-login-page o feature/carrito-compra`

**fix/**: Para corregir un bug que no es urgente en producción. Ejemplo: `fix/error-calculo-iva`

**hotfix/**: Para corregir un bug urgente y crítico directamente en producción. Ejemplo: `hotfix/caida-servidor-pagos`

**release/**: Para preparar una nueva versión antes de pasarla a producción (ajustes finales, testing). Ejemplo: `release/v2.3.0`

**chore/**: Para tareas de mantenimiento, configuración o limpieza que no afectan a funcionalidad. Ejemplo: `chore/actualizar-dependencias`

**refactor/**: Para reestructurar código sin cambiar su comportamiento. Ejemplo: `refactor/optimizar-servicio-usuarios`

**docs/**: Para cambios solo en documentación. Ejemplo: `docs/actualizar-readme`

**test/**: Para añadir o modificar tests. Ejemplo: `test/cobertura-modulo-pagos`

**style/**: Para cambios de formato o estilo de código. Ejemplo: `style/aplicar-eslint`

**ci/**: Para cambios en pipelines de integración continua. Ejemplo: `ci/mejorar-workflow-deploy`