import{i as Y}from"./vue-demi@0.13.11_vue@3.2.47-71ba0ef2.js";import{n as H,k as q,m as I,h as Z,a as L,d as B,t as G,g as $,o as A,q as T}from"./@vue_reactivity@3.2.47-7641be53.js";import{g as tt,m as et,i as st,n as nt,f as ot}from"./@vue_runtime-core@3.2.47-cef9f3df.js";/*!
  * pinia v2.0.33
  * (c) 2023 Eduardo San Martin Morote
  * @license MIT
  */let D;const R=t=>D=t,J=Symbol();function C(t){return t&&typeof t=="object"&&Object.prototype.toString.call(t)==="[object Object]"&&typeof t.toJSON!="function"}var k;(function(t){t.direct="direct",t.patchObject="patch object",t.patchFunction="patch function"})(k||(k={}));function ht(){const t=H(!0),o=t.run(()=>q({}));let s=[],e=[];const r=I({install(u){R(r),r._a=u,u.provide(J,r),u.config.globalProperties.$pinia=r,e.forEach(a=>s.push(a)),e=[]},use(u){return!this._a&&!Y?e.push(u):s.push(u),this},_p:s,_a:null,_e:t,_s:new Map,state:o});return r}const N=()=>{};function V(t,o,s,e=N){t.push(o);const r=()=>{const u=t.indexOf(o);u>-1&&(t.splice(u,1),e())};return!s&&$()&&A(r),r}function g(t,...o){t.slice().forEach(s=>{s(...o)})}function x(t,o){t instanceof Map&&o instanceof Map&&o.forEach((s,e)=>t.set(e,s)),t instanceof Set&&o instanceof Set&&o.forEach(t.add,t);for(const s in o){if(!o.hasOwnProperty(s))continue;const e=o[s],r=t[s];C(r)&&C(e)&&t.hasOwnProperty(s)&&!L(e)&&!B(e)?t[s]=x(r,e):t[s]=e}return t}const ct=Symbol();function rt(t){return!C(t)||!t.hasOwnProperty(ct)}const{assign:y}=Object;function ut(t){return!!(L(t)&&t.effect)}function ft(t,o,s,e){const{state:r,actions:u,getters:a}=o,f=s.state.value[t];let j;function b(){f||(s.state.value[t]=r?r():{});const v=T(s.state.value[t]);return y(v,u,Object.keys(a||{}).reduce((m,d)=>(m[d]=I(ot(()=>{R(s);const p=s._s.get(t);return a[d].call(p,p)})),m),{}))}return j=W(t,b,o,s,e,!0),j}function W(t,o,s={},e,r,u){let a;const f=y({actions:{}},s),j={deep:!0};let b,v,m=I([]),d=I([]),p;const _=e.state.value[t];!u&&!_&&(e.state.value[t]={}),q({});let O;function F(c){let n;b=v=!1,typeof c=="function"?(c(e.state.value[t]),n={type:k.patchFunction,storeId:t,events:p}):(x(e.state.value[t],c),n={type:k.patchObject,payload:c,storeId:t,events:p});const h=O=Symbol();nt().then(()=>{O===h&&(b=!0)}),v=!0,g(m,n,e.state.value[t])}const z=u?function(){const{state:n}=s,h=n?n():{};this.$patch(S=>{y(S,h)})}:N;function K(){a.stop(),m=[],d=[],e._s.delete(t)}function M(c,n){return function(){R(e);const h=Array.from(arguments),S=[],w=[];function U(i){S.push(i)}function X(i){w.push(i)}g(d,{args:h,name:c,store:l,after:U,onError:X});let E;try{E=n.apply(this&&this.$id===t?this:l,h)}catch(i){throw g(w,i),i}return E instanceof Promise?E.then(i=>(g(S,i),i)).catch(i=>(g(w,i),Promise.reject(i))):(g(S,E),E)}}const Q={_p:e,$id:t,$onAction:V.bind(null,d),$patch:F,$reset:z,$subscribe(c,n={}){const h=V(m,c,n.detached,()=>S()),S=a.run(()=>et(()=>e.state.value[t],w=>{(n.flush==="sync"?v:b)&&c({storeId:t,type:k.direct,events:p},w)},y({},j,n)));return h},$dispose:K},l=Z(Q);e._s.set(t,l);const P=e._e.run(()=>(a=H(),a.run(()=>o())));for(const c in P){const n=P[c];if(L(n)&&!ut(n)||B(n))u||(_&&rt(n)&&(L(n)?n.value=_[c]:x(n,_[c])),e.state.value[t][c]=n);else if(typeof n=="function"){const h=M(c,n);P[c]=h,f.actions[c]=n}}return y(l,P),y(G(l),P),Object.defineProperty(l,"$state",{get:()=>e.state.value[t],set:c=>{F(n=>{y(n,c)})}}),e._p.forEach(c=>{y(l,a.run(()=>c({store:l,app:e._a,pinia:e,options:f})))}),_&&u&&s.hydrate&&s.hydrate(l.$state,_),b=!0,v=!0,l}function bt(t,o,s){let e,r;const u=typeof o=="function";typeof t=="string"?(e=t,r=u?s:o):(r=t,e=t.id);function a(f,j){const b=tt();return f=f||b&&st(J,null),f&&R(f),f=D,f._s.has(e)||(u?W(e,o,r,f):ft(e,r,f)),f._s.get(e)}return a.$id=e,a}export{ht as c,bt as d};
