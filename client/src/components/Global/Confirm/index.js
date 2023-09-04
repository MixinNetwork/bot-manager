import _confirm from './Confirm.vue'

const confirm = {}
confirm.install = vue => {
  const ConfirmCon = vue.extend(_confirm)
  const ins = new ConfirmCon()
  ins.$mount(document.createElement('div'))
  document.body.appendChild(ins.$el)
  vue.prototype.$confirm = (optionOrTitle, onSuccess) => {
    if (typeof optionOrTitle === "string") {
      ins.title = optionOrTitle
      ins.onSuccess = async () => {
        let t = await onSuccess()
        if (t !== false) {
          ins.visible = false
        }
      }
    } else if (typeof optionOrTitle === "object") {
      const { confirm, cancel, success } = optionOrTitle
      if (confirm) ins.confirm = confirm
      if (cancel) ins.cancel = cancel
      if (success) ins.onSuccess = async () => {
        let t = await success()
        if (t !== false) {
          ins.visible = false
        }
      }
    }
    ins.visible = true
  }
}
export default confirm
