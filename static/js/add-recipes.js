function addIngredientField() {
    const container = document.getElementById("ingredients");
    const newChild = document.createElement("div");
    const id = container.childElementCount;
    newChild.id = `ingredient${id}`;

    const labelName = document.createElement("label");
    labelName.innerText = "name:";
    newChild.appendChild(labelName);

    const inputName = document.createElement("input");
    inputName.type = "text";
    inputName.name = `Ingredients[${id}].Name`;
    newChild.appendChild(inputName);

    const labelCount = document.createElement("label");
    labelCount.innerText = "count:";
    newChild.appendChild(labelCount);

    const inputCount = document.createElement("input");
    inputCount.type = "text";
    inputCount.name = `Ingredients[${id}].Count`;
    newChild.appendChild(inputCount);

    const buttonRemove = document.createElement("a");
    buttonRemove.onclick = () => { removeIngredientField(newChild.id) };
    buttonRemove.innerText = "X";
    newChild.appendChild(buttonRemove);

    container.appendChild(newChild);
}

function removeIngredientField(id) {
    const container = document.getElementById("ingredients");
    const field = document.getElementById(id);
    const inputs = field.getElementsByTagName("input");
    const inputName = inputs.item(0);
    const inputCount = inputs.item(1);
    const lastField = document.getElementById(`ingredient${container.childElementCount-1}`);
    container.removeChild(field);
    if (id !== lastField.id) {
        lastField.id = id;
        const lastInputs = lastField.getElementsByTagName("input");
        const lastName = lastInputs.item(0);
        const lastCount = lastInputs.item(1);
        lastName.id = inputName.id;
        lastName.name = inputName.name;
        lastCount.id = inputCount.id;
        lastCount.name = inputCount.name
    }
}
