<template>
  <Modal :show="actionState">
    <div class="modal-container" v-if="activeGrower.avatar_url">
      <div class="head">
        <img :src="activeGrower.avatar_url">
        <span>{{activeGrower.full_name}}</span>
        <i>Mixin ID : {{activeGrower.identity_number}}</i>
      </div>
      <div class="text">
        <label>博客</label>
        <a :target="activeGrower.blog && '_blank'"
           :href="activeGrower.blog ? activeGrower.blog : 'javascript:void(0)'"
        >{{activeGrower.blog}}</a>
      </div>
      <div class="text">
        <label>Twitter</label>
        <a :target="activeGrower.twitter && '_blank'"
           :href="activeGrower.twitter ? activeGrower.twitter.replace('@','https://twitter.com/') : 'javascript:void(0)'"
        >{{activeGrower.twitter}}</a>
      </div>
      <div class="text">
        <label>GitHub</label>
        <a :target="activeGrower.github && '_blank'"
           :href="activeGrower.github ? activeGrower.github : 'javascript:void(0)'"
        >{{activeGrower.github}}</a>
      </div>
      <div class="text">
        <label>审核意见</label>

        <textarea v-model="auditOpinion"/>
      </div>

      <div class="btns">
        <button @click="clickBtn('cancel')" class="btn cancel">取消</button>
        <button @click="clickBtn('pass')" class="btn pass">通过</button>
        <button @click="clickBtn('reject')" class="btn reject">拒绝</button>
      </div>
    </div>
  </Modal>
</template>

<script>
  import Modal from '@/components/Modal'
  import { mapState, mapActions } from 'vuex'

  export default {
    name: "GrowthAudit",
    components: { Modal },
    data() {
      return {
        auditOpinion: ''
      }
    },
    watch: {
      activeGrower(v) {
        this.auditOpinion = v.auditopinion
      }
    },
    computed: {
      ...mapState('growth', ['actionState', 'activeGrower'])
    },
    methods: {
      ...mapActions('growth', ['updateState', 'toggleState']),
      clickAction(grower) {
        this.activeGrower = grower
      },
      closeModal() {
        this.$DC('growth', { actionState: false })
      },
      async clickBtn(option) {
        if (option === 'cancel') return closeAction.call(this)
        this.$DC('growth', { auditOpinion: this.auditOpinion })
        let t = await this.updateState(option)
        if (t === true) {
          this.toggleState()
          return closeAction.call(this)
        }
      }
    },
  }

  function closeAction() {
    this.$DC('growth', { actionState: false })
  }
</script>

<style lang="scss" scoped>

  a {
    color: #4C4471;
  }

  .modal-container {
    width: 700px;
    height: 600px;

    .head {
      text-align: center;
      display: flex;
      flex-direction: column;
      align-items: center;

      span {
        color: #000;
        font-size: 18px;
      }

      i {
        color: #D8D8D8;
      }
    }


    img {
      width: 72px;
      height: 72px;
      border-radius: 50%;
      margin: 13px 0 23px 0;
    }
  }


  .text {
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 26px 0;

    label {
      width: 100px;
      color: #000;
      font-weight: bold;
    }

    a {
      flex: 1;
      line-height: 60px;
      height: 60px;
      max-width: 490px;
      overflow: hidden;
      text-overflow: ellipsis;
      background: #F6F9FF;
      border-radius: 4px;
      padding: 0 20px;
      font-size: 14px;
    }

    textarea {
      flex: 1;
      max-width: 490px;
      background: #F6F9FF;
      height: 68px;
      border-radius: 4px;
      padding: 20px;

      resize: none;
    }
  }

  .btns {
    display: flex;
    justify-content: center;

    .btn {
      width: 140px;
      height: 40px;
      border-radius: 4px;
      margin-right: 30px;
      cursor: pointer;

      &:last-child {
        margin-right: 0;
      }
    }

    .cancel {
      background-color: #E5E7EC;
    }

    .pass {
      background-color: #1BAC2C;
      color: #fff;
    }

    .reject {
      background-color: #A50A29;
      color: #fff;
    }
  }
</style>