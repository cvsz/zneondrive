const data={
  districts:[
    ["Central Arcology","NOVA Authority","Civic towers, official events and surveillance.",["story","legal race","social"]],
    ["Glass Crown","VANTEX","Corporate skyline, luxury contracts and DRIVE ZERO power.",["story","luxury jobs","social"]],
    ["Foundry 9","Mechanist Guild","Fabrication, scrapyards, Garage 17 and early build progression.",["build","salvage","story"]],
    ["Docklands","Independent","Freight, towing, convoy routes and contraband stories.",["delivery","convoy","story"]],
    ["Neon Mile","Iron Wolves","Nightlife, car meets and street-racing celebrity.",["street race","social","nightlife"]],
    ["Dust Verge","Free Roads","Desert routes, abandoned energy sites and off-road play.",["off-road","exploration","convoy"]],
    ["Greenbelt","NOVA Authority","Reclaimed ecology, protected roads and companion utility.",["exploration","pet utility","time attack"]],
    ["Old Grid","Contested","Pre-war networks, archives and DRIVE ZERO archaeology.",["story","archives","salvage"]],
    ["Skyway","NOVA Authority","Elevated expressways, championship and time attack.",["circuit","time attack","championship"]],
    ["Black Circuit","Iron Wolves","Hidden race network unlocked by reputation.",["underground race","high risk","story"]]
  ],
  chapters:[
    ["01","The Broken Machine","Build the first vehicle and become mobile."],
    ["02","Street Reputation","Work the city and earn a name."],
    ["03","Underground","Choose who to trust in the Black Circuit."],
    ["04","The Championship","Qualify, compete and become visible."],
    ["05","VANTEX","Turn the prototype into evidence."],
    ["06","Drive Zero","Survive pursuit and understand the control layer."],
    ["07","War for Nova","Choose the future governance of mobility."]
  ],
  quests:[
    ["MQ001","Arrival After Midnight"],["MQ002","Garage 17"],["MQ003","A Chassis With No Name"],
    ["MQ004","Scrap Rights"],["MQ005","Borrowed Tools"],["MQ006","Dead Grid"],
    ["MQ007","First Fuel"],["MQ008","Legacy ECU"],["MQ009","Missing Fastener"],
    ["MQ010","Cold Crank"],["MQ011","First Ignition"],["MQ012","Roadworthy"]
  ],
  roles:{
    Street:{desc:"Balanced daily build for city jobs and mixed routes.",stats:{Speed:62,Grip:62,Control:72,Torque:58,Utility:60}},
    Drift:{desc:"High steering angle and rotation; trades stability for style/control.",stats:{Speed:58,Grip:38,Control:88,Torque:66,Utility:42}},
    Circuit:{desc:"Grip, braking and consistency for sanctioned track racing.",stats:{Speed:76,Grip:88,Control:80,Torque:64,Utility:28}},
    Drag:{desc:"Launch and straight-line acceleration with narrow specialization.",stats:{Speed:90,Grip:52,Control:42,Torque:92,Utility:20}},
    "Off-road":{desc:"Travel, suspension and recovery capability for Dust Verge.",stats:{Speed:52,Grip:70,Control:74,Torque:84,Utility:90}},
    Delivery:{desc:"Reliable transport setup with cargo utility and efficiency.",stats:{Speed:50,Grip:64,Control:70,Torque:62,Utility:96}}
  }
};

const i18n={
  th:{navWorld:"โลก",navStory:"เนื้อเรื่อง",navVehicle:"รถ",navSlice:"Vertical Slice",navArchitecture:"ระบบ",navRoadmap:"Roadmap",eyebrow:"ONLINE OPEN WORLD RPG · VEHICLE BUILDING · RACING · SOCIAL",heroLead:"เกมออนไลน์ที่ “รถ” เป็นตัวตนที่สองของผู้เล่น — สร้าง ซ่อม แข่ง และสะสมประวัติของรถคันเดียวไปจนเป็นตำนานของ Server",startDemo:"เริ่ม Demo",seeArchitecture:"ดูสถาปัตยกรรม",metricQuests:"Main Quests",metricDistricts:"Districts",metricCharacters:"Characters",metricFactions:"Factions",metricChapters:"Story Chapters",worldTitle:"หนึ่งเมือง หลายชีวิต หลายอำนาจ",worldText:"เมืองปี 2097 ที่ถนน เทเลเมทรี และระบบยานยนต์กลายเป็นโครงสร้างพื้นฐานเชิงอำนาจ",storyTitle:"จากรถพัง → สู่สงครามชิงอนาคตของ NOVA CITY",twistText:"รถคันแรกไม่ใช่เศษเหล็กธรรมดา แต่คือ Prototype ที่ ECU ซ่อนกุญแจของระบบที่สามารถควบคุมยานพาหนะทั้งเมืองได้",vehicleTitle:"รถหนึ่งคัน เปลี่ยนบทบาทได้ทั้งเกม",vehicleText:"เลือก Build Role เพื่อดูแนวทาง trade-off ของรถคันเดิม — ไม่มี VIP boost ต่อ performance",fairnessTitle:"VIP เพิ่มพื้นที่ ไม่เพิ่มพลัง",garageSlot:"Garage Slot",garageSlots:"Garage Slots",fair1:"Performance signature เท่ากันเมื่อใช้ Build เดียวกัน",fair2:"ไม่มี hidden horsepower / grip / durability",fair3:"Ranked progression เปิดครบสำหรับ Free Player",sliceTitle:"Garage 17 → First Ignition",sliceText:"กดเดินเรื่องเพื่อจำลอง proof-of-concept ของ first-session loop และ authoritative state",nextQuest:"Next Quest",reset:"Reset",architectureTitle:"Server-authoritative by design",architectureText:"Client ส่ง intent — server เป็นผู้ตัดสิน state ที่มีมูลค่า แข่งขัน หรือคงอยู่ระยะยาว",roadmapTitle:"ขายเป็น Milestone ที่พิสูจน์ได้ ไม่ขายคำว่า “เสร็จ”",estimateNote:"* Planning range only; final schedule depends on art fidelity, target platforms, concurrency and team size.",decisionTitle:"พรุ่งนี้ต้องตัดสินใจ 4 เรื่อง",decision1:"PC first หรือ PC + Console?",decision2:"Visual fidelity ระดับไหน?",decision3:"Concurrency milestone แรกเท่าไร?",decision4:"Vertical Slice / Alpha / Full Program?"},
  en:{navWorld:"World",navStory:"Story",navVehicle:"Vehicle",navSlice:"Vertical Slice",navArchitecture:"Architecture",navRoadmap:"Roadmap",eyebrow:"ONLINE OPEN WORLD RPG · VEHICLE BUILDING · RACING · SOCIAL",heroLead:"A persistent online game where the vehicle becomes the player's second identity — built, repaired, raced and remembered until one machine becomes a server legend.",startDemo:"Start Demo",seeArchitecture:"See Architecture",metricQuests:"Main Quests",metricDistricts:"Districts",metricCharacters:"Characters",metricFactions:"Factions",metricChapters:"Story Chapters",worldTitle:"One city. Many lives. Competing powers.",worldText:"In 2097 roads, telemetry and vehicle networks have become strategic infrastructure.",storyTitle:"From a broken machine to a war for NOVA CITY",twistText:"The starter car is not random scrap. It is a DRIVE ZERO prototype whose ECU hides an artifact capable of changing control over city-wide mobility.",vehicleTitle:"One vehicle, multiple roles",vehicleText:"Switch build roles to show meaningful trade-offs on the same vehicle — with no VIP performance boost.",fairnessTitle:"VIP adds capacity, not power",garageSlot:"Garage Slot",garageSlots:"Garage Slots",fair1:"Identical builds produce identical competitive signatures",fair2:"No hidden horsepower / grip / durability",fair3:"Full ranked progression remains available to free players",sliceTitle:"Garage 17 → First Ignition",sliceText:"Advance the story to simulate the first-session proof and its authoritative state.",nextQuest:"Next Quest",reset:"Reset",architectureTitle:"Server-authoritative by design",architectureText:"Clients send intent; authoritative servers decide persistent, valuable and competitive state.",roadmapTitle:"Sell evidence-backed milestones, not the word “finished”",estimateNote:"* Planning range only; final schedule depends on art fidelity, target platforms, concurrency and team size.",decisionTitle:"Four decisions to leave the meeting with",decision1:"PC first or PC + Console?",decision2:"What visual-fidelity target?",decision3:"What first concurrency milestone?",decision4:"Vertical Slice / Alpha / Full Program?"}
};

let lang="th", currentQuest=0, currentRole="Street";
const qs=s=>document.querySelector(s);
const qsa=s=>[...document.querySelectorAll(s)];

function applyLanguage(){
  document.documentElement.lang=lang;
  qsa("[data-i18n]").forEach(el=>{const key=el.dataset.i18n;if(i18n[lang][key])el.textContent=i18n[lang][key]});
  qs("#langToggle").textContent=lang==="th"?"EN":"TH";
}
qs("#langToggle").addEventListener("click",()=>{lang=lang==="th"?"en":"th";applyLanguage()});

function renderDistricts(){
  const grid=qs("#districtGrid");
  data.districts.forEach((d,i)=>{
    const b=document.createElement("button");b.type="button";b.className="district";
    b.innerHTML=`<span>${String(i+1).padStart(2,"0")} · ${d[1]}</span><strong>${d[0]}</strong>`;
    b.addEventListener("click",()=>selectDistrict(i));grid.appendChild(b);
  });
  selectDistrict(2);
}
function selectDistrict(i){
  qsa(".district").forEach((e,x)=>e.classList.toggle("active",x===i));
  const d=data.districts[i];
  qs("#districtDetail").innerHTML=`<p class="eyebrow">${d[1]}</p><h3>${d[0]}</h3><p>${d[2]}</p><div class="activity-list">${d[3].map(x=>`<span>${x}</span>`).join("")}</div>`;
}
function renderChapters(){
  qs("#chapterTimeline").innerHTML=data.chapters.map(c=>`<article class="chapter"><b>CHAPTER ${c[0]}</b><strong>${c[1]}</strong><small>${c[2]}</small></article>`).join("");
}
function renderRoles(){
  const box=qs("#roleButtons");
  Object.keys(data.roles).forEach(role=>{
    const b=document.createElement("button");b.type="button";b.className="role-btn";b.textContent=role;
    b.addEventListener("click",()=>selectRole(role));box.appendChild(b);
  });
  selectRole(currentRole);
}
function selectRole(role){
  currentRole=role;qsa(".role-btn").forEach(b=>b.classList.toggle("active",b.textContent===role));
  const r=data.roles[role];qs("#buildRole").textContent=role.toUpperCase();qs("#buildDescription").textContent=r.desc;
  qs("#statBars").innerHTML=Object.entries(r.stats).map(([k,v])=>`<div class="stat"><span>${k}</span><div class="bar"><i style="width:${v}%"></i></div><b>${v}</b></div>`).join("");
}
function renderQuests(){
  qs("#questStepper").innerHTML=data.quests.map((q,i)=>`<div class="quest ${i<currentQuest?"done":i===currentQuest?"current":""}"><b>${q[0]}</b><small>${q[1]}</small><span class="status">${i<currentQuest?"✓":i===currentQuest?"●":"·"}</span></div>`).join("");
  updateState();
}
function updateState(){
  const q=data.quests[currentQuest][0];
  qs("#stateQuest").textContent=q;
  qs("#stateVehicle").textContent=currentQuest>=2?"ND-PROTOTYPE-001":"UNCLAIMED";
  qs("#stateBuild").textContent=currentQuest>=9?"3":currentQuest>=4?"2":currentQuest>=2?"1":"0";
  qs("#stateIgnition").textContent=currentQuest>=10?"ENABLED":"LOCKED";
  qs("#stateRoadworthy").textContent=currentQuest>=11?"YES":"NO";
  const notes=[
    "Message acknowledged. Objective activation is idempotent.",
    "Garage checkpoint committed; Maya relationship initialized.",
    "Starter vehicle bound to the primary character.",
    "Salvage grant protected from duplicate retries.",
    "Build revision appended; stale expected revisions are rejected.",
    "Traversal objective persists across reconnect.",
    "Authored choice committed once.",
    "Legacy ECU discovery cannot fabricate reward state.",
    "Recovered item grant is idempotent.",
    "Ignition readiness validates the active build.",
    "First Ignition references the exact active build revision.",
    "Chapter boundary, reward and Roadworthy state commit idempotently."
  ];
  qs("#stateNote").textContent=notes[currentQuest];
}
qs("#nextQuest").addEventListener("click",()=>{if(currentQuest<data.quests.length-1)currentQuest++;renderQuests()});
qs("#prevQuest").addEventListener("click",()=>{if(currentQuest>0)currentQuest--;renderQuests()});
qs("#resetQuest").addEventListener("click",()=>{currentQuest=0;renderQuests()});

renderDistricts();renderChapters();renderRoles();renderQuests();applyLanguage();
