(window.webpackJsonp=window.webpackJsonp||[]).push([["chunk-76fb"],{D8vC:function(e,n,t){},GV1Q:function(e,n,t){"use strict";var o=t("D8vC");t.n(o).a},c11S:function(e,n,t){"use strict";var o=t("gTgX");t.n(o).a},gTgX:function(e,n,t){},ntYl:function(e,n,t){"use strict";t.r(n);var o={name:"Login",components:{LangSelect:t("ETGp").a},data:function(){return{loginForm:{username:"",password:""},loginRules:{username:[{required:!0,trigger:"blur",validator:function(e,n,t){n.length<6?t(new Error("Please enter the correct user name")):t()}}],password:[{required:!0,trigger:"blur",validator:function(e,n,t){n.length<6?t(new Error("The password can not be less than 6 digits")):t()}}]},loginDisable:!0,passwordType:"password",pwdIcon:"eye",loading:!1,showDialog:!1,redirect:void 0}},computed:{loginBtnDisable:function(){return this.loginForm.password.length<4||this.loginForm.username.length<4}},watch:{$route:{handler:function(e){this.redirect=e.query&&e.query.redirect},immediate:!0}},methods:{showPwd:function(){"password"===this.passwordType?(this.passwordType="",this.pwdIcon="eye_open"):(this.passwordType="password",this.pwdIcon="eye")},handleLogin:function(){var e=this;this.$refs.loginForm.validate(function(n){if(!n)return console.log("error submit!!"),!1;e.loading=!0,e.$store.dispatch("Login",e.loginForm).then(function(){e.loading=!1,e.$router.push({path:e.redirect||"/"})}).catch(function(){e.loading=!1})})}}},s=(t("c11S"),t("GV1Q"),t("KHd+")),i=Object(s.a)(o,function(){var e=this,n=e.$createElement,t=e._self._c||n;return t("div",{staticClass:"login-container"},[t("el-form",{ref:"loginForm",staticClass:"login-form",attrs:{model:e.loginForm,"auto-complete":"on","label-position":"left"}},[t("div",{staticClass:"title-container"},[t("lang-select",{staticClass:"set-language"}),e._v(" "),t("h3",{staticClass:"title"},[e._v(e._s(e.$t("login.title")))])],1),e._v(" "),t("el-form-item",{attrs:{prop:"username"}},[t("span",{staticClass:"svg-container"},[t("svg-icon",{attrs:{"icon-class":"user"}})],1),e._v(" "),t("el-input",{attrs:{placeholder:e.$t("login.username"),name:"username",type:"text","auto-complete":"on"},model:{value:e.loginForm.username,callback:function(n){e.$set(e.loginForm,"username",n)},expression:"loginForm.username"}})],1),e._v(" "),t("el-form-item",{attrs:{prop:"password"}},[t("span",{staticClass:"svg-container"},[t("svg-icon",{attrs:{"icon-class":"password"}})],1),e._v(" "),t("el-input",{attrs:{type:e.passwordType,placeholder:e.$t("login.password"),"auto-complete":"on"},nativeOn:{keyup:function(n){return"button"in n||!e._k(n.keyCode,"enter",13,n.key,"Enter")?e.handleLogin(n):null}},model:{value:e.loginForm.password,callback:function(n){e.$set(e.loginForm,"password",n)},expression:"loginForm.password"}}),e._v(" "),t("span",{staticClass:"show-pwd",on:{click:e.showPwd}},[t("svg-icon",{attrs:{"icon-class":e.pwdIcon}})],1)],1),e._v(" "),t("el-button",{staticStyle:{width:"100%","margin-bottom":"30px"},attrs:{loading:e.loading,disabled:e.loginBtnDisable,type:"primary"},nativeOn:{click:function(n){return n.preventDefault(),e.handleLogin(n)}}},[e._v("\n      "+e._s(e.$t("login.logIn"))+"\n    ")]),e._v(" "),t("el-form-item",[t("div",[t("span",{staticStyle:{color:"#eef1f6"}},[e._v("测试账号：admin")]),e._v(" "),t("span",{staticStyle:{color:"#eef1f6"}},[e._v("密码：admin")])])])],1)],1)},[],!1,null,"6b19f46a",null);i.options.__file="index.vue";n.default=i.exports}}]);