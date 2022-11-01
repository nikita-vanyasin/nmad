<template>
  <div class="map-container">
    <GoogleMap
        :api-key="apiKey"
        style="width: 100%; height: calc(100vh - 75px)"
        :zoom="2"
        :zoomControl="true"
    >
      <CustomMarker v-for="nomadLocation in nomadLocations" :options="nomadLocation">
        <div class="marker">
          <span style="font-size: 1.125rem">
            <a class="link" :href="nomadLocation.profile_url" target="_blank" title="send message">{{nomadLocation.username}}</a>
          </span>
        </div>
      </CustomMarker>
    </GoogleMap>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { GoogleMap, CustomMarker } from 'vue3-google-map'
import {Api} from "@/api";

export default defineComponent({
  components: { GoogleMap, CustomMarker },
  data() {
    return {
      nomadLocations: [],
      apiKey: import.meta.env.VITE_GOOGLE_API_KEY
    }
  },
  methods: {
    async fetchNomadLocations() {
      let chatId = this.$route.params.chat_id
      const resp = await Api.get(`/list/${chatId}`)

      this.nomadLocations = resp.data.map((loc) => {
        return {
          position: {
            lat: loc.lat,
            lng: loc.lng
          },
          anchorPoint: 'BOTTOM_CENTER',
          username: loc.username,
          profile_url: loc.profile_url,
        }
      })
    }
  },
  async mounted() {
    await this.fetchNomadLocations()
  }
})
</script>
<style scoped>
.marker {
  text-align: center;
  background-color: white;
  border-radius: 15px;
  padding: 4px 8px;
}

.marker .link,
.marker .link:visited {
  text-decoration: none;
  color: #5e5e5e;
  opacity: 0.85;
}
</style>