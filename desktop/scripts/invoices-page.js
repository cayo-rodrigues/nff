export function createInvoicesPage() {
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
            <div class="invoices-form__sections-container">
                ${invoicesFormSection(1)}
            </div>

            <button class="invoices-form__button" type="submit">Emitir Notas Fiscais</button>
            <button class="invoices-form__add-section-button" type="button">+</button>
        </form>
    `

    const sectionsContainer = document.querySelector('.invoices-form__sections-container')
    sectionsContainer.addEventListener('click', ({ target }) => {
        if (target.id && target.id.includes('open-dialog-button')) {
            document.getElementById(`items-dialog-${target.dataset.invoiceId}`).showModal()
        }
        else if (target.id && target.id.includes('close-dialog-button')) {
            document.getElementById(`items-dialog-${target.dataset.invoiceId}`).close()
        }
        else if (target.id && target.id.includes('add-item-button')) {
            const invoiceId = target.dataset.invoiceId
            const dialogSectionsContainer = document.getElementById(`dialog-sections-container-${invoiceId}`)
            const sectionId = dialogSectionsContainer.childElementCount + 1
            dialogSectionsContainer.innerHTML += invoiceItemsSection(invoiceId, sectionId)
        }
    })

    const addSectionButton = contentCore.querySelector('.invoices-form__add-section-button')
    addSectionButton.addEventListener('click', () => {
        const invoiceId = sectionsContainer.childElementCount + 1
        sectionsContainer.innerHTML += invoicesFormSection(invoiceId)
    })

    const form = contentCore.querySelector('.invoices-form')
    form.addEventListener('submit', submitInvoicesForm)
}


function invoicesFormSection(id) {
    return `
        <section class="invoices-form__section">
            <h3>Nota Fiscal ${id}</h3>
            <div id="${id}" class="invoices-form__inputs-container">

                <div class="invoices-form__input">
                    <label for="sender-input-${id}">Remetente</label>
                    <select name="sender" id="sender-input-${id}">
                        <option value="sender-id">Emerson</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="recipient-input-${id}">Destinatário</label>
                    <select name="recipient" id="recipient-input-${id}">
                        <option value="emerson-id">Emerson</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="operation-input-${id}">Natureza da Operação</label>
                    <select name="operation" id="operation-input-${id}">
                        <option value="VENDA">VENDA</option>
                        <option value="REMESSA">REMESSA</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="gta-input-${id}">GTA</label>
                    <input type="gta" name="gta" id="gta-input-${id}">
                </div>
                <div class="invoices-form__input">
                    <label for="cfop-input-${id}">CFOP</label>
                    <select name="cfop" id="cfop-input-${id}">
                        <option value="5101">5101</option>
                        <option value="5102">5102</option>
                        <option value="5103">5103</option>
                        <option value="5105">5105</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="shipping-input-${id}">Frete</label>
                    <input type="number" step=0.01 name="shipping" id="shipping-input-${id}">
                </div>
                <div class="invoices-form__input">
                    <label for="add_shipping_to_total_value-input-${id}">Adicionar Frete ao Total</label>
                    <select name="add_shipping_to_total_value" id="add_shipping_to_total_value-input-${id}">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="is_final_customer-input-${id}">Consumidor Final</label>
                    <select name="is_final_customer" id="is_final_customer-input-${id}">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="icms-input-${id}">Contribuinte ICMS</label>
                    <select name="icms" id="icms-input-${id}">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="custom_file_name-input-${id}">Nome do Arquivo</label>
                    <input type="text" name="custom_file_name" id="custom_file_name-input-${id}">
                </div>
                <div class="invoices-form__input">
                    <label for="extra_notes-input-${id}">Informações Complementares</label>
                    <input type="text" name="extra_notes" id="extra_notes-input-${id}">
                </div>
                <div class="invoices-form__input">
                    <label>Itens da Nota Fiscal</label>
                    <button 
                        type="button"
                        id="open-dialog-button-${id}"
                        class="invoices-form__items-dialog-button"
                        data-invoice-id="${id}"
                    >
                        Gerenciar Itens
                    </button>
                </div>
                
                ${manageInvoiceItemsDialog(id)}

            </div>
        </section>
    `
}

function manageInvoiceItemsDialog(invoiceId) {
    return `
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
                ${invoiceItemsSection(invoiceId, 1)}
            </div>

        </dialog>
    `
}

function invoiceItemsSection(invoiceId, sectionId) {
    return `
        <section class="invoices-form__items-section">
            <h4>Item ${sectionId}</h4>

            <div id="${invoiceId}-${sectionId}" class="invoices-form__inputs-container">
                
                <div class="invoices-form__input">
                    <label for="group-input-${invoiceId}-${sectionId}">Grupo</label>
                    <select name="group" id="group-input-${invoiceId}-${sectionId}" data-section="items">
                        <option value="Gado bovino para corte">Gado bovino para corte</option>
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
                        <option value="Nacional">Nacional</option>
                    </select>
                </div>

                <div class="invoices-form__input">
                    <label for="unity_of_measurement-input-${invoiceId}-${sectionId}">Unidade de medida</label>
                    <select name="unity_of_measurement" id="unity_of_measurement-input-${invoiceId}-${sectionId}" data-section="items">
                        <option value="CB">CB</option>
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
        </section>

    `
}

export function cancelInvoicesPage() {
    const contentCore = document.querySelector("#content__core")
    contentCore.innerHTML = ""
    contentCore.innerHTML = `
        Cancel invoices!
    `
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
            // if this is true, then it means that in this loop begins the next form section
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
