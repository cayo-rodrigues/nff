import { toTitleCase, svgIcons, listsDataToOptionTags } from "./helpers.js"

export async function listEntitiesPage() {
    document.querySelector("#current-tab-title").innerText = "Entidades"

    const className = "sub-menu__item--selected"

    const subEntry1 = document.querySelector("#sub-entry-1")
    subEntry1.innerText = "Listar Entidades"
    subEntry1.classList.add(className)

    const subEntry2 = document.querySelector("#sub-entry-2")
    subEntry2.innerText = "Registrar Entidades"
    subEntry2.classList.remove(className)

    const contentCore = document.querySelector("#content__core")
    contentCore.innerHTML = ""

    const entities = await pywebview.api.get_entities()
    const listsData = await pywebview.api.get_lists_data(
        'user_type_options, street_type_options'
    )
    const optionTags = listsDataToOptionTags(listsData)

    const entitiesList = document.createElement("ul")
    entitiesList.className = "entities-list"

    for (let entity of entities) {
        entitiesList.append(makeEntityCard(entity, optionTags))
    }

    entitiesList.addEventListener('click', handleEntityActions)
    contentCore.append(entitiesList)
    contentCore.querySelectorAll('.entity-form').forEach((form) => {
        form.addEventListener('submit', (event) => submitEntityForm(event, 'update'))
    })
}

function makeEntityCard(entity, optionTags) {
    const entityCard = document.createElement('li')
    entityCard.className = 'entities-list__card'
    entityCard.id = entity.id
    entityCard.innerHTML = `
      <div class="entity-card__shallow-info">
          <h4 class="shallow-info__data">${entity.name}</h4>
          <div class="shallow-info__data">
            <b>CPF/CNPJ</b>: <span>${entity.cpf_cnpj || '-'}</span>
          </div>
          <div class="shallow-info__data">
            <b>IE</b>: <span>${entity.ie || '-'}</span>
          </div>
          <div class="shallow-info__data">
            <b>Bairro</b>: <span>${entity.neighborhood || '-'}</span>
          </div>
      </div>

      <div class="entity-card__action-icons">
        <span class="entity-card__update-icon">
            ${svgIcons.update}
        </span>
        <span class="entity-card__delete-icon">
            ${svgIcons.delete}
        </span>
      </div>
    `
    entityCard.append(entityDetailsDialog(entity), updateEntityDialog(entity, optionTags))
    return entityCard
}

function updateEntityDialog(entity, optionTags) {
    const dialog = document.createElement('dialog')
    dialog.className = "entity-card__update-dialog"
    dialog.id = `update-dialog-${entity.id}`
    dialog.dataset.entityId = entity.id
    dialog.innerHTML = `
        <heading class="update-entity__heading">
            <h3>Editando ${entity.name}</h3>
            <button class="update-dialog__close-button">X</button>
        </heading>
        
        <form class="entity-form" data-entity-id="${entity.id}">
            <div class="entity-form__inputs-container">

                <div class="entity-form__input">
                    <label for="name-input">Nome</label>
                    <input type="text" name="name" id="name-input" value="${entity.name}">
                </div>
                <div class="entity-form__input">
                    <label for="email-input">Email</label>
                    <input type="email" name="email" id="email-input" value="${entity.email}">
                </div>
                <div class="entity-form__input">
                    <label for="user_type-input">Tipo</label>
                    <select name="user_type" id="user_type-input">
                        ${optionTags.user_type_options}
                    </select>
                </div>
                <div class="entity-form__input">
                    <label for="cpf_cnpj-input">CPF/CNPJ</label>
                    <input type="text" name="cpf_cnpj" id="cpf_cnpj-input" value="${entity.cpf_cnpj}">
                </div>
                <div class="entity-form__input">
                    <label for="ie-input">Inscrição Estadual</label>
                    <input type="text" name="ie" id="ie-input" value="${entity.ie}">
                </div>
                <div class="entity-form__input">
                    <label for="password-input">Senha</label>
                    <input type="password" name="password" id="password-input" value="${entity.password}">
                </div>
                <div class="entity-form__input">
                    <label for="postal_code-input">CEP</label>
                    <input type="text" name="postal_code" id="postal_code-input" value="${entity.postal_code}">
                </div>
                <div class="entity-form__input">
                    <label for="neighborhood-input">Bairro</label>
                    <input type="text" name="neighborhood" id="neighborhood-input" value="${entity.neighborhood}">
                </div>
                <div class="entity-form__input">
                    <label for="street_type-input">Logradouro (tipo)</label>
                    <select name="street_type" id="street_type-input">
                        ${optionTags.street_type_options}
                    </select>
                </div>
                <div class="entity-form__input">
                    <label for="street_name-input">Logradouro (nome)</label>
                    <input type="text" name="street_name" id="street_name-input" value="${entity.street_name}">
                </div>
                <div class="entity-form__input">
                    <label for="number-input">Número</label>
                    <input type="text" name="number" id="number-input" value="${entity.number}">
                </div>

            </div>

            <button class="entity-form__button" type="submit">Atualizar Entidade</button>
        </form>
    `

    dialog.querySelector(`option[value="${entity.user_type}"]`).selected = true
    dialog.querySelector(`option[value="${entity.street_type}"]`).selected = true

    return dialog
}

function entityDetailsDialog(entity) {
    const dialog = document.createElement('dialog')
    dialog.className = 'entity-card__details-dialog'
    dialog.id = `details-dialog-${entity.id}`
    dialog.dataset.entityId = entity.id
    dialog.innerHTML = `
          <heading class="entity-details__heading">
            <h3>${entity.name}</h3>
            <button class="details-dialog__close-button">X</button>
          </heading>
          <hr/>
          <section class="details-dialog__info-section">
              <ul>
                  <li class="entity-details__item">
                      <b>Tipo</b>: <span>${entity.user_type}</span>
                  </li>
                  <li class="entity-details__item">
                      <b>CPF/CNPJ</b>: <span>${entity.cpf_cnpj}</span>
                  </li>
                  <li class="entity-details__item">
                      <b>IE</b>: <span>${entity.ie}</span>
                  </li>
                  <li class="entity-details__item">
                      <b>Email</b>: <span>${entity.email}</span>
                  </li>
                  <li class="entity-details__item">
                      <b>Senha</b>: <span>${entity.password}</span>
                  </li>
              </ul>
              <ul>
                  <li class="entity-details__item">
                      <b>CEP</b>: <span>${entity.postal_code}</span>
                  </li>
                  <li class="entity-details__item">
                      <b>Bairro</b>: <span>${entity.neighborhood}</span>
                  </li>
                  <li class="entity-details__item">
                      <b>Logradouro (tipo)</b>: <span>${entity.street_type}</span>
                  </li>
                  <li class="entity-details__item">
                      <b>Logradouro (nome)</b>: <span>${entity.street_name}</span>
                  </li>
                  <li class="entity-details__item">
                      <b>Número</b>: <span>${entity.number}</span>
                  </li>
              </ul>
          </section>
    `
    return dialog
}

async function handleEntityActions(event) {
    if (event.target.matches('span.entity-card__delete-icon, span.entity-card__delete-icon *')) {
        const entityCard = event.target.closest('li.entities-list__card')
        await pywebview.api.delete_entity(entityCard.id)
        entityCard.remove()
    }
    else if (event.target.matches('span.entity-card__update-icon, span.entity-card__update-icon *')) {
        document.getElementById(`update-dialog-${event.target.closest('li.entities-list__card').id}`).showModal()
    }
    else if (event.target.className === 'update-dialog__close-button') {
        document.getElementById(`update-dialog-${event.target.closest('li.entities-list__card').id}`).close()
    }
    else if (event.target.matches('div.entity-card__shallow-info, div.entity-card__shallow-info *')) {
        document.getElementById(`details-dialog-${event.target.closest('li.entities-list__card').id}`).showModal()
    }
    else if (event.target.className === 'details-dialog__close-button') {
        document.getElementById(`details-dialog-${event.target.closest('li.entities-list__card').id}`).close()
    }
}

export async function createEntitiesPage() {
    const listsData = await pywebview.api.get_lists_data(
        'user_type_options, street_type_options'
    )
    const options = listsDataToOptionTags(listsData)
    const contentCore = document.querySelector("#content__core")
    contentCore.innerHTML = `
        <form class="entity-form">
            <div class="entity-form__inputs-container">

                <div class="entity-form__input">
                    <label for="name-input">Nome</label>
                    <input type="text" name="name" id="name-input">
                </div>
                <div class="entity-form__input">
                    <label for="email-input">Email</label>
                    <input type="email" name="email" id="email-input">
                </div>
                <div class="entity-form__input">
                    <label for="user_type-input">Tipo</label>
                    <select name="user_type" id="user_type-input">
                        ${options.user_type_options}
                    </select>
                </div>
                <div class="entity-form__input">
                    <label for="cpf_cnpj-input">CPF/CNPJ</label>
                    <input type="text" name="cpf_cnpj" id="cpf_cnpj-input">
                </div>
                <div class="entity-form__input">
                    <label for="ie-input">Inscrição Estadual</label>
                    <input type="text" name="ie" id="ie-input">
                </div>
                <div class="entity-form__input">
                    <label for="password-input">Senha</label>
                    <input type="password" name="password" id="password-input">
                </div>
                <div class="entity-form__input">
                    <label for="postal_code-input">CEP</label>
                    <input type="text" name="postal_code" id="postal_code-input">
                </div>
                <div class="entity-form__input">
                    <label for="neighborhood-input">Bairro</label>
                    <input type="text" name="neighborhood" id="neighborhood-input">
                </div>
                <div class="entity-form__input">
                    <label for="street_type-input">Logradouro (tipo)</label>
                    <select name="street_type" id="street_type-input">
                        ${options.street_type_options}
                    </select>
                </div>
                <div class="entity-form__input">
                    <label for="street_name-input">Logradouro (nome)</label>
                    <input type="text" name="street_name" id="street_name-input">
                </div>
                <div class="entity-form__input">
                    <label for="number-input">Número</label>
                    <input type="text" name="number" id="number-input">
                </div>

            </div>

            <button class="entity-form__button" type="submit">Registrar Entidade</button>
        </form>
    `
    const form = contentCore.querySelector('.entity-form')
    form.addEventListener('submit', (event) => submitEntityForm(event, 'create'))
}

async function submitEntityForm(event, action) {
    event.preventDefault()
    const form = event.target

    const formData = {}

    for (const child of form) {
        if (child.name) {
            if (["name", "neighborhood", "street_name"].includes(child.name)) {
                child.value = toTitleCase(child.value)
            }
            formData[child.name] = child.value
        }
    }

    if (action === 'create') {
        await pywebview.api.register_entity(formData)
        listEntitiesPage()
    }
    else if (action === 'update') {
        const entityId = form.dataset.entityId
        await pywebview.api.update_entity(formData, entityId)
        const entityCard = document.getElementById(entityId)
        const newEntityData = Object.assign({ id: entityId }, formData)
        entityCard.parentNode.replaceChild(makeEntityCard(newEntityData), entityCard)
    }
}
