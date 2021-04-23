<template>
  <div class="input-meta">
    <el-form :model="internalValue">
      <el-form-item inline="true" v-for="(kv,index) in internalValue.kvs" :key="index" :label="kv.key" class="meta-pairs">
        <el-input placeholder="Value" v-model="kv.value" v-on:change="onChange(kv)" class="meta-value"/>
        <el-button type="danger" icon="el-icon-delete" v-on:click="onDelete(kv)" class="meta-action-btn" />
      </el-form-item>
    </el-form>

    <el-form ref="newMetaForm" :model="newKv" :rules="newMetaRule">
      <el-form-item inline="true" class="meta-pairs" prop="key">
        <el-input placeholder="Key" v-model="newKv.key" required="true" inline-message="true" class="meta-key"/>
        <el-input placeholder="Value" v-model="newKv.value" class="meta-value" />
        <el-button icon="el-icon-plus"  v-on:click="onAdd(newKv)" class="meta-action-btn"/>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>

function pairsFromObject(kv){
  let ret = [];
  for (const [key, value] of Object.entries(kv)) {
    ret.push({key, value})
  }

  return {
    kvs: ret
  };
}


function setPairsToValue(vm){
  let newKeys = [];
  vm.internalValue.kvs.forEach(function(kv){
    vm.value[kv.key] = kv.value
    newKeys.push(kv.key)
  })

  let deleteKeys = Object.keys(vm.value).filter(v => !newKeys.includes(v))
  deleteKeys.forEach(deleteKey => delete vm.value[deleteKey])
}

function validateDuplicateKey(vm, newKey){
  if(newKey === ""){
    return
  }

  for (let i = 0, l = vm.internalValue.kvs.length; i < l; i++){
    let kv = vm.internalValue.kvs[i]
    if (kv.key !== newKey){
      continue
    }
    return new Error('Key 已存在');
  }
}

export default {
  name: "MetaInput",
  props: {
    value: {
      type: Object,
      default: function() {
        return {}
      }
    },
  },

  data() {
    return {
      internalValue: this.value ? pairsFromObject(this.value) : [],
      newKv: {key: "", value: ""},
      newMetaRule: {
        key:[
            {
              required: true,
              message: 'key 必须填写',
              trigger: 'blur',
            },
            {
              required: true,
              trigger: 'change',
              message: 'key 已存在',
              vm: this,
              validator: function (rule, value, callback){
                callback(validateDuplicateKey(rule.vm, value))
              }
            }
          ]
      },
    }
  },

  methods: {
    onDelete: function(kv){
      let self = this
      this.internalValue.kvs.forEach(function(_kv, i){
          if(_kv.key === kv.key){
            self.internalValue.kvs.splice(i, 1)
            return
          }
      });
    },
    onAdd: function (kv){
      if(kv.key === "" || kv.value === ""){
        return
      }

      let self = this
      this.$refs["newMetaForm"].validate(function (valid){
        if(!valid){
          return;
        }
        self.internalValue.kvs.push({key: kv.key, value: kv.value})
        kv.key = ""
        kv.value = "";
      });

    },

    onChange: function (kv){
      this.internalValue.kvs.forEach(function(_kv){
        if(_kv.key === kv.key){
          _kv.value = kv.value
          return
        }
      });
    },

  },
  watch: {
    internalValue: {
      deep: true,
      handler: function (){
        window.se = this
        setPairsToValue(this)
      }
    },
  },

}
</script>

<style scoped>
  .input-meta{
    position: relative;
    padding-left: 10px;
  }

  .meta-value {
    display: inline-block;
    width: 80%;
    min-width: 350px;
  }

  .meta-pairs .meta-key {
    display: inline-block;
    width: 10%;
    min-width: 100px;
    margin-right: 1em;
  }


</style>