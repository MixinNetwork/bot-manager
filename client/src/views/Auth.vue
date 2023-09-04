<template>
  <Loading />
</template>

<script>
import Loading from "@/components/Loading";
import { mapActions } from "vuex";
export default {
  name: "Auth",
  components: { Loading },
  methods: {
    ...mapActions("user", ["authenticate"])
  },
  async mounted() {
    const code = getUrlParameter("code");
    const return_to = getUrlParameter("return_to") || "/";
    let status = await this.authenticate(code);
    if (status) this.$router.push(return_to);
    else this.$bus.toAuth();
  }
};

function getUrlParameter(name) {
  name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
  var regex = new RegExp("[\\?&]" + name + "=([^&#]*)");
  var results = regex.exec(window.location.search);
  return results === null
    ? ""
    : decodeURIComponent(results[1].replace(/\+/g, " "));
}
</script>