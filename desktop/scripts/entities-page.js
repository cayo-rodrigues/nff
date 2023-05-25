export function listEntitiesPage() {
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

  // MOCK
  const entities = [
    {
      name: "Emerson Frivasa",
      type: "Produtor Rural",
      ie: "123456789",
      cpf_cnpj: "987654321",
      email: "aushuahsuahs@gmail.com",
      password: "123456",
      zip_code: "37508000",
      neighborhood: "Laje",
      street_type: "Rua",
      street_name: "Olegário Maciel",
      number: "47",
    },
  ]
  // MOCK

  const entitiesList = document.createElement("ul")
  entitiesList.className = "entities-list"

  for (let entity of entities) {
    entitiesList.innerHTML += `
      <li class="entities-list__card">

        <div class="entity-card__heading" >
          <h3 class="entity-card__name">${entity.name}</h3>
          <div class="entity-form__input" class="entity-card__action-icons">
            <span>Edit Icon</span>
            <span>Delete Icon</span>
          </div>
        </div>
        
        <div class="entity-card__contents" >
          <ul class="entity-card__section">
            <li><b>Tipo:</b> <span>${entity.type}</span></li>
            <li><b>IE:</b> <span>${entity.ie}</span></li>
            <li><b>CPF/CNPJ:</b> <span>${entity.cpf_cnpj}</span></li>
            <li><b>E-mail:</b> <span>${entity.email}</span></li>
            <li><b>Senha:</b> <span>${entity.password}</span></li>
          </ul>
          
          <ul class="entity-card__section">
            <li><b>CEP:</b> <span>${entity.zip_code}</span></li>
            <li><b>Bairro:</b> <span>${entity.neighborhood}</span></li>
            <li><b>Logradouro (tipo):</b> <span>${entity.street_type}</span></li>
            <li><b>Logradouro (nome):</b> <span>${entity.street_name}</span></li>
            <li><b>Número:</b> <span>${entity.number}</span></li>
          </ul>
        </div>

      </li>
    `
  }

  contentCore.append(entitiesList)
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
                    <label for="type-input">Tipo</label>
                    <select name="entity_type" id="type-input">
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
                    <input type="text" name="address_number" id="number-input">
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
            formData[child.name] = child.value
        }
    }
    
    const response = await pywebview.api.register_entity(formData)
    console.log(response)
}
