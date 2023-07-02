import { toTitleCase } from "./helpers.js"

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

    const entitiesList = document.createElement("ul")
    entitiesList.className = "entities-list"

    for (let entity of entities) {
        entitiesList.innerHTML += `
      <li class="entities-list__card" id="${entity.id}">

          <div class="entity-card__shallow-info">
              <h4 class="entity-card__name shallow-info__data">${entity.name}</h4>
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

          ${entityDetailsDialog(entity)}

          <div class="entity-card__action-icons">
            <span class="entity-card__update-icon">
              <svg viewBox="0 0 1024 1024" fill="#000000" class="icon" version="1.1"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"><path d="M574.4 590.4l-3.2 7.2 1.6 8L608 740.8l8 33.6 28-20L760 672l5.6-4 2.4-6.4 220-556.8 8.8-22.4-22.4-8.8-140-55.2-21.6-8-8.8 20.8-229.6 559.2z m244-528l140 55.2-13.6-30.4-220 556.8 8-10.4-116 82.4 36 13.6-33.6-135.2-0.8 15.2 229.6-560-29.6 12.8z" fill=""></path><path d="M872 301.6l-107.2-40c-7.2-2.4-10.4-10.4-8-17.6l8-20.8c2.4-7.2 10.4-10.4 17.6-8l107.2 40c7.2 2.4 10.4 10.4 8 17.6l-8 20.8c-2.4 7.2-10.4 10.4-17.6 8zM718.4 645.6l-107.2-40c-7.2-2.4-10.4-10.4-8-17.6l8-20.8c2.4-7.2 10.4-10.4 17.6-8l107.2 40c7.2 2.4 10.4 10.4 8 17.6l-8 20.8c-2.4 7.2-10.4 10.4-17.6 8zM900.8 224l-107.2-40c-7.2-2.4-10.4-10.4-8-17.6l8-20.8c2.4-7.2 10.4-10.4 17.6-8l107.2 40c7.2 2.4 10.4 10.4 8 17.6l-8 20.8c-2.4 7.2-10.4 11.2-17.6 8z" fill=""></path><path d="M930.4 965.6H80c-31.2 0-56-24.8-56-56V290.4c0-31.2 24.8-56 56-56h576c13.6 0 24 10.4 24 24s-10.4 24-24 24H80c-4 0-8 4-8 8v619.2c0 4 4 8 8 8h850.4c4 0 8-4 8-8V320c0-13.6 10.4-24 24-24s24 10.4 24 24v589.6c0 31.2-24.8 56-56 56z" fill=""></path><path d="M366.4 490.4H201.6c-13.6 0-25.6-11.2-25.6-25.6 0-13.6 11.2-25.6 25.6-25.6h165.6c13.6 0 25.6 11.2 25.6 25.6-0.8 14.4-12 25.6-26.4 25.6zM409.6 584h-208c-13.6 0-25.6-11.2-25.6-25.6 0-13.6 11.2-25.6 25.6-25.6h208c13.6 0 25.6 11.2 25.6 25.6-0.8 14.4-12 25.6-25.6 25.6zM441.6 676.8h-240c-13.6 0-25.6-11.2-25.6-25.6 0-13.6 11.2-25.6 25.6-25.6h240c13.6 0 25.6 11.2 25.6 25.6-0.8 14.4-12 25.6-25.6 25.6z" fill=""></path></g></svg>
            </span>
            <span class="entity-card__delete-icon">
                <svg fill="#000000" viewBox="0 0 24 24"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"><path d="M5.755,20.283,4,8H20L18.245,20.283A2,2,0,0,1,16.265,22H7.735A2,2,0,0,1,5.755,20.283ZM21,4H16V3a1,1,0,0,0-1-1H9A1,1,0,0,0,8,3V4H3A1,1,0,0,0,3,6H21a1,1,0,0,0,0-2Z"></path></g></svg>
            </span>
          </div>
      </li>
    `
    }

    entitiesList.addEventListener('click', handleEntityActions)
    contentCore.append(entitiesList)
}

function entityDetailsDialog(entity) {
    return `
      <dialog class="entity-card__details-dialog" id="details-dialog-${entity.id}" data-entity-id="${entity.id}">
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
      </dialog>
    `
}

async function handleEntityActions(event) {
    if (event.target.matches('span.entity-card__delete-icon, span.entity-card__delete-icon *')) {
        const entityCard = event.target.closest('li.entities-list__card')
        await pywebview.api.delete_entity(entityCard.id)
        entityCard.remove()
    } else if (event.target.matches('span.entity-card__update-icon, span.entity-card__update-icon *')) {
        console.log('WANTS TO UPDATE ENTITY ID =>', event.target.closest('li.entities-list__card').id)
    } else if (event.target.matches('div.entity-card__shallow-info, div.entity-card__shallow-info *')) {
        document.getElementById(`details-dialog-${event.target.closest('li.entities-list__card').id}`).showModal()
    } else if (event.target.className === 'details-dialog__close-button') {
        document.getElementById(`details-dialog-${event.target.closest('li.entities-list__card').id}`).close()
    }
}

export function createEntitiesPage() {
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
                        <option value="Produtor Rural">Produtor Rural</option>
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
                        <option value="Rua">Rua</option>
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
    form.addEventListener('submit', submitEntityForm)
}

async function submitEntityForm(event) {
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

    const response = await pywebview.api.register_entity(formData)
    console.log(response)
}
