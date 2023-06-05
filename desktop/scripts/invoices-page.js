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
    })

    const addSectionButton = contentCore.querySelector('.invoices-form__add-section-button')
    addSectionButton.addEventListener('click', () => {
        const invoiceId = sectionsContainer.childElementCount + 1
        sectionsContainer.innerHTML += invoicesFormSection(invoiceId)
    })
}


function invoicesFormSection(id) {
    return `
        <section class="invoices-form__section">
            <h3>Nota Fiscal ${id}</h3>
            <div id="${id}" class="invoices-form__inputs-container">

                <div class="invoices-form__input">
                    <label for="sender-input">Remetente</label>
                    <select name="sender" id="sender-input">
                        <option value="sender-id">Emerson</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="recipient-input">Destinatário</label>
                    <select name="recipient" id="recipient-input">
                        <option value="emerson-id">Emerson</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="operation-input">Natureza da Operação</label>
                    <select name="operation" id="operation-input">
                        <option value="VENDA">VENDA</option>
                        <option value="REMESSA">REMESSA</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="gta-input">GTA</label>
                    <input type="gta" name="gta" id="gta-input">
                </div>
                <div class="invoices-form__input">
                    <label for="cfop-input">CFOP</label>
                    <select name="cfop" id="cfop-input">
                        <option value="5101">5101</option>
                        <option value="5102">5102</option>
                        <option value="5103">5103</option>
                        <option value="5105">5105</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="shipping-input">Frete</label>
                    <input type="number" step=0.01 name="shipping" id="shipping-input">
                </div>
                <div class="invoices-form__input">
                    <label for="add_shipping_to_total_value-input">Adicionar Frete ao Total</label>
                    <select name="add_shipping_to_total_value" id="add_shipping_to_total_value-input">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="is_final_customer-input">Consumidor Final</label>
                    <select name="is_final_customer" id="is_final_customer-input">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="icms-input">Contribuinte ICMS</label>
                    <select name="icms" id="icms-input">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="custom_file_name-input">Nome do Arquivo</label>
                    <input type="text" name="custom_file_name" id="custom_file_name-input">
                </div>
                <div class="invoices-form__input">
                    <label for="extra_notes-input">Informações Complementares</label>
                    <input type="text" name="extra_notes" id="extra_notes-input">
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

function manageInvoiceItemsDialog(id) {
    return `
        <dialog id="items-dialog-${id}" data-invoice-id="${id}" class="invoice-items-dialog">
            <div class="invoice-items-dialog__heading">
                <h3>Itens da Nota Fiscal ${id}</h3>
                    
                <div class="invoice-items-dialog__buttons-container">
                    <button 
                        type="button"
                        class="invoice-items-dialog__button invoice-items-dialog__confirm-button"
                    >
                        Confirmar
                    </button>
                    <button
                        type="button"
                        class="invoice-items-dialog__button invoice-items-dialog__cancel-button"
                        id="close-dialog-button-${id}"
                        data-invoice-id="${id}"
                    >
                        Cancelar
                    </button>
                </div>
            </div>

            <hr/>
            
            <div class="invoices-form__inputs-container">
                
                <div class="invoices-form__input">
                    <label for="group-input">Grupo</label>
                    <select name="group" id="group-input">
                        <option value="Gado bovino para corte">Gado bovino para corte</option>
                    </select>
                </div>

                <div class="invoices-form__input">
                    <label for="ncm-input">NCM</label>
                    <input type="text" name="ncm" id="ncm-input">
                </div>

                <div class="invoices-form__input">
                    <label for="description-input">Descrição</label>
                    <input type="text" name="description" id="description-input">
                </div>

                <div class="invoices-form__input">
                    <label for="origin-input">Origem</label>
                    <select name="origin" id="origin-input">
                        <option value="Nacional">Nacional</option>
                    </select>
                </div>

                <div class="invoices-form__input">
                    <label for="unity_of_measurement-input">Unidade de medida</label>
                    <select name="unity_of_measurement" id="unity_of_measurement-input">
                        <option value="CB">CB</option>
                    </select>
                </div>

                <div class="invoices-form__input">
                    <label for="quantity-input">Quantidade</label>
                    <input type="number" step=0.01 name="quantity" id="quantity-input">
                </div>

                <div class="invoices-form__input">
                    <label for="value_per_unity-input">Valor Unitário</label>
                    <input type="number" step=0.01 name="value_per_unity" id="value_per_unity-input">
                </div>

            </div>

        </dialog>
    `
}

export function cancelInvoicesPage() {
    const contentCore = document.querySelector("#content__core")
    contentCore.innerHTML = ""
    contentCore.innerHTML = `
        Cancel invoices!
    `
}
