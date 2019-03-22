<script>
import { Line, mixins } from 'vue-chartjs'
import 'chartjs-plugin-streaming'

export default {
  extends: Line,
  mixins: [mixins.reactiveProp],
  props: ['options'],
  mounted() {
    this.renderChart( this.chartData, this.options )
  },
}
</script>