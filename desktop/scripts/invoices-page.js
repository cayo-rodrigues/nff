import { entitiesToOptionTags, listsDataToOptionTags } from "./helpers.js"

export async function createInvoicesPage() {
    document.querySelector("#current-tab-title").innerText = "Notas Fiscais"

    const className = "sub-menu__item--selected"

    const subEntry1 = document.querySelector("#sub-entry-1")
    subEntry1.innerText = "Emitir Notas Fiscais"
    subEntry1.classList.add(className)

    const subEntry2 = document.querySelector("#sub-entry-2")
    subEntry2.innerText = "Cancelar Notas Fiscais"
    subEntry2.classList.remove(className)

    const contentCore = document.querySelector("#content__core")
    contentCore.innerHTML = `
        <form class="invoices-form">
            <div class="invoices-form__sections-container"></div>

            <button class="invoices-form__button" type="submit">Emitir Notas Fiscais</button>
            <button class="invoices-form__add-section-button" type="button">+</button>
        </form>
    `

    const sectionsContainer = document.querySelector('.invoices-form__sections-container')
    sectionsContainer.addEventListener('click', ({ target }) => {
        if (target.id && target.id.includes('open-dialog-button')) {
            document.getElementById(`items-dialog-${target.dataset.invoiceId}`).showModal()
        }
        else if (target.id && (target.id.includes('close-dialog-button') || target.id.includes('confirm-items-button'))) {
            // CONFIRM BUTTON SHOULD KEEP THE ITEMS
            // CANCEL BUTTON SHOULD ERASE THE ITEMS
            document.getElementById(`items-dialog-${target.dataset.invoiceId}`).close()
        }
        else if (target.id && target.id.includes('add-item-button')) {
            const invoiceId = target.dataset.invoiceId
            const dialogSectionsContainer = document.getElementById(`dialog-sections-container-${invoiceId}`)
            const sectionId = dialogSectionsContainer.childElementCount + 1

            newInvoiceItemSection(invoiceId, sectionId, dialogSectionsContainer, optionTags)
        }
    })

    const entities = await pywebview.api.get_entities()
    const listsData = await pywebview.api.get_lists_data(
        'operation_options, cfop_options, icms_options, group_options, ' +
        'origin_options, unity_of_measurement_options, boolean_options'
    )
    const optionTags = {
        entities: entitiesToOptionTags(entities),
    }
    Object.assign(optionTags, listsDataToOptionTags(listsData))

    const addSectionButton = contentCore.querySelector('.invoices-form__add-section-button')
    addSectionButton.addEventListener('click', () => {
        const invoiceId = sectionsContainer.childElementCount + 1
        newInvoiceSection(invoiceId, sectionsContainer, optionTags)
    })

    const form = contentCore.querySelector('.invoices-form')
    form.addEventListener('submit', submitInvoicesForm)

    newInvoiceSection(1, sectionsContainer, optionTags)
}


function newInvoiceSection(invoiceId, sectionsContainer, optionTags) {
    const newFormSection = document.createElement('section')
    newFormSection.className = "invoices-form__section"
    newFormSection.innerHTML = `
        <h3>Nota Fiscal ${invoiceId}</h3>
        <div id="${invoiceId}" class="invoices-form__inputs-container">

            <div class="invoices-form__input">
                <label for="sender-input-${invoiceId}">Remetente</label>
                <select name="sender" id="sender-input-${invoiceId}">
                    ${optionTags.entities}
                </select>
            </div>
            <div class="invoices-form__input">
                <label for="recipient-input-${invoiceId}">Destinatário</label>
                <select name="recipient" id="recipient-input-${invoiceId}">
                    ${optionTags.entities}
                </select>
            </div>
            <div class="invoices-form__input">
                <label for="operation-input-${invoiceId}">Natureza da Operação</label>
                <select name="operation" id="operation-input-${invoiceId}">
                    ${optionTags.operation_options}
                </select>
            </div>
            <div class="invoices-form__input">
                <label for="gta-input-${invoiceId}">GTA</label>
                <input type="text" name="gta" id="gta-input-${invoiceId}">
            </div>
            <div class="invoices-form__input">
                <label for="cfop-input-${invoiceId}">CFOP</label>
                <select name="cfop" id="cfop-input-${invoiceId}">
                    ${optionTags.cfop_options}
                </select>
            </div>
            <div class="invoices-form__input">
                <label for="shipping-input-${invoiceId}">Frete</label>
                <input type="number" step=0.01 name="shipping" id="shipping-input-${invoiceId}">
            </div>
            <div class="invoices-form__input">
                <label for="add_shipping_to_total_value-input-${invoiceId}">Adicionar Frete ao Total</label>
                <select name="add_shipping_to_total_value" id="add_shipping_to_total_value-input-${invoiceId}">
                    ${optionTags.boolean_options}
                </select>
            </div>
            <div class="invoices-form__input">
                <label for="is_final_customer-input-${invoiceId}">Consumidor Final</label>
                <select name="is_final_customer" id="is_final_customer-input-${invoiceId}">
                    ${optionTags.boolean_options}
                </select>
            </div>
            <div class="invoices-form__input">
                <label for="icms-input-${invoiceId}">Contribuinte ICMS</label>
                <select name="icms" id="icms-input-${invoiceId}">
                    ${optionTags.icms_options}
                </select>
            </div>
            <div class="invoices-form__input">
                <label for="custom_file_name-input-${invoiceId}">Nome do Arquivo</label>
                <input type="text" name="custom_file_name" id="custom_file_name-input-${invoiceId}">
            </div>
            <div class="invoices-form__input">
                <label for="extra_notes-input-${invoiceId}">Informações Complementares</label>
                <input type="text" name="extra_notes" id="extra_notes-input-${invoiceId}">
            </div>
            <div class="invoices-form__input">
                <label>Itens da Nota Fiscal</label>
                <button 
                    type="button"
                    id="open-dialog-button-${invoiceId}"
                    class="invoices-form__items-dialog-button"
                    data-invoice-id="${invoiceId}"
                >
                    Gerenciar Itens
                </button>
            </div>
            
            <dialog id="items-dialog-${invoiceId}" class="invoice-items-dialog">
                <div class="invoice-items-dialog__heading">
                    <h3>Itens da Nota Fiscal ${invoiceId}</h3>
                        
                    <div class="invoice-items-dialog__buttons-container">
                        <button
                            type="button"
                            class="invoice-items-dialog-button"
                            id="add-item-button-${invoiceId}"
                            data-invoice-id="${invoiceId}"
                        >
                            +
                        </button>
                        <button 
                            type="button"
                            class="invoice-items-dialog__button invoice-items-dialog__confirm-button"
                            id="confirm-items-button-${invoiceId}"
                            data-invoice-id="${invoiceId}"
                        >
                            Confirmar
                        </button>
                        <button
                            type="button"
                            class="invoice-items-dialog__button invoice-items-dialog__cancel-button"
                            id="close-dialog-button-${invoiceId}"
                            data-invoice-id="${invoiceId}"
                        >
                            Cancelar
                        </button>
                    </div>
                </div>

                <hr/>
                
                <div id="dialog-sections-container-${invoiceId}" class="invoices-form__dialog-sections-container">
                </div>

            </dialog>

        </div>
    `

    sectionsContainer.append(newFormSection)
    newInvoiceItemSection(invoiceId, 1, document.getElementById(`dialog-sections-container-${invoiceId}`), optionTags)
}

function newInvoiceItemSection(invoiceId, sectionId, dialogSectionsContainer, optionTags) {
    const newItemSection = document.createElement('section')
    newItemSection.className = "invoices-form__items-section"
    newItemSection.innerHTML = `
        <h4>Item ${sectionId}</h4>

        <div id="${invoiceId}-${sectionId}" class="invoices-form__inputs-container">
            
            <div class="invoices-form__input">
                <label for="group-input-${invoiceId}-${sectionId}">Grupo</label>
                <select name="group" id="group-input-${invoiceId}-${sectionId}" data-section="items">
                    ${optionTags.group_options}
                </select>
            </div>

            <div class="invoices-form__input">
                <label for="ncm-input-${invoiceId}-${sectionId}">NCM</label>
                <input type="text" name="ncm" id="ncm-input-${invoiceId}-${sectionId}" data-section="items">
            </div>

            <div class="invoices-form__input">
                <label for="description-input-${invoiceId}-${sectionId}">Descrição</label>
                <input type="text" name="description" id="description-input-${invoiceId}-${sectionId}" data-section="items">
            </div>

            <div class="invoices-form__input">
                <label for="origin-input-${invoiceId}-${sectionId}">Origem</label>
                <select name="origin" id="origin-input-${invoiceId}-${sectionId}" data-section="items">
                    ${optionTags.origin_options}
                </select>
            </div>

            <div class="invoices-form__input">
                <label for="unity_of_measurement-input-${invoiceId}-${sectionId}">Unidade de medida</label>
                <select name="unity_of_measurement" id="unity_of_measurement-input-${invoiceId}-${sectionId}" data-section="items">
                    ${optionTags.unity_of_measurement_options}
                </select>
            </div>

            <div class="invoices-form__input">
                <label for="quantity-input-${invoiceId}-${sectionId}">Quantidade</label>
                <input type="number" step=0.01 name="quantity" id="quantity-input-${invoiceId}-${sectionId}" data-section="items">
            </div>

            <div class="invoices-form__input">
                <label for="value_per_unity-input-${invoiceId}-${sectionId}">Valor Unitário</label>
                <input type="number" step=0.01 name="value_per_unity" id="value_per_unity-input-${invoiceId}-${sectionId}" data-section="items">
            </div>

        </div>
    `

    dialogSectionsContainer.append(newItemSection)
}

export async function cancelInvoicesPage() {
    const contentCore = document.querySelector("#content__core")
    contentCore.innerHTML = ""
    contentCore.innerHTML = `
        <form class="invoices-form">
            <div class="invoices-form__sections-container"></div>

            <button class="invoices-form__button" type="submit">Cancelar Notas Fiscais</button>
            <button class="invoices-form__add-section-button" type="button">+</button>
        </form>
    `

    const entities = await pywebview.api.get_entities()
    const optionTags = {
        entities: entitiesToOptionTags(entities)
    }

    const sectionsContainer = document.querySelector('.invoices-form__sections-container')

    const addSectionButton = contentCore.querySelector('.invoices-form__add-section-button')
    addSectionButton.addEventListener('click', () => {
        const cancelingId = sectionsContainer.childElementCount + 1
        newCancelingSection(cancelingId, sectionsContainer, optionTags)
    })

    const form = contentCore.querySelector('.invoices-form')
    form.addEventListener('submit', submitCancelingsForm)

    newCancelingSection(1, sectionsContainer, optionTags)
}

function newCancelingSection(cancelingId, sectionsContainer, optionTags) {
    const newFormSection = document.createElement('section')
    newFormSection.className = "invoices-form__section"
    newFormSection.innerHTML = `
        <h3>Cancelamento ${cancelingId}</h3>
        <div id="${cancelingId}" class="invoices-form__inputs-container">

            <div class="invoices-form__input">
                <label for="entity-${cancelingId}">Entidade</label>
                <select name="entity" id="entity-${cancelingId}">
                    ${optionTags.entities}
                </select>
            </div>
            <div class="invoices-form__input">
                <label for="invoice-id-${cancelingId}">Número da nota</label>
                <input type="text" name="invoice_id" id="invoice-id-${cancelingId}">
            </div>
            <div class="invoices-form__input">
                <label for="year-${cancelingId}">Ano</label>
                <input type="text" name="year" id="year-${cancelingId}">
            </div>
            <div class="invoices-form__input">
                <label for="justification-${cancelingId}">Justificativa</label>
                <input type="text" name="justification" id="justification-${cancelingId}">
            </div>
            
        </div>
    `

    sectionsContainer.append(newFormSection)
}

async function submitInvoicesForm(event) {
    event.preventDefault()
    const form = event.target

    let invoiceSectionData = {
        items: []
    }
    let itemSectionData = {}
    const invoicesData = []

    for (const child of form) {
        if (child.name) {
            // if this is true, then it means that the next invoice section begins now
            if (invoiceSectionData.hasOwnProperty(child.name)) {
                invoiceSectionData.items.push(Object.assign({}, itemSectionData))
                invoicesData.push(Object.assign({}, invoiceSectionData))
                invoiceSectionData = {
                    items: []
                }
                itemSectionData = {}
            }
            // here it means that this is the next item section within an invoice section
            if (itemSectionData.hasOwnProperty(child.name)) {
                invoiceSectionData.items.push(Object.assign({}, itemSectionData))
                itemSectionData = {}
            }
            child.dataset.section === 'items'
                ? itemSectionData[child.name] = child.value
                : invoiceSectionData[child.name] = child.value
        }
    }
    // applied only to the last invoice section 
    invoicesData.push(invoiceSectionData)
    invoiceSectionData.items.push(itemSectionData)


    const response = await pywebview.api.create_invoices(invoicesData)
    console.log(response)
}

async function submitCancelingsForm(event) {
    event.preventDefault()
    const form = event.target

    const cancelingsData = []
    let sectionData = {}

    for (const child of form) {
        if (child.name) {
            if (sectionData.hasOwnProperty(child.name)) {
                cancelingsData.push(Object.assign({}, sectionData))
                sectionData = {}
            }
            sectionData[child.name] = child.value
        }
    }
    cancelingsData.push(sectionData)

    const response = await pywebview.api.cancel_invoices(cancelingsData)
    console.log(response)
}
