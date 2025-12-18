# Go-Study2 å‰ç«¯ï¼ˆNext.js 14ï¼‰

åŸºäº Next.js 14ï¼ˆApp Routerï¼‰+ TypeScript + Ant Designï¼Œé™æ€å¯¼å‡ºç”±åç«¯æ‰˜ç®¡ã€‚

## å¿«é€Ÿå¼€å§‹

```bash
cd frontend
npm install
npm run dev   # http://localhost:3000
```

ç¯å¢ƒå˜é‡ï¼š

- `NEXT_PUBLIC_API_URL`ï¼šåç«¯ API åŸºå€ï¼ˆé»˜è®¤ `/api/v1`ï¼Œç”Ÿäº§å»ºè®®è®¾ä¸ºå®Œæ•´åŸŸåï¼‰ã€‚
- `NEXT_PUBLIC_API_URL`ï¼šåç«¯ API åŸºå€ã€‚å»ºè®®åœ¨æœ¬åœ°å¼€å‘æ—¶æŒ‡å‘åç«¯æœåŠ¡ï¼ˆä¾‹å¦‚ `http://localhost:8080`ï¼‰ã€‚

å¼€å‘æç¤ºï¼š
- å¦‚æœæœªè®¾ç½® `NEXT_PUBLIC_API_URL`ï¼ŒNext.js å¼€å‘æœåŠ¡å™¨ä¼šå°è¯•å°†ä»¥ `/api/v1` å¼€å¤´çš„è¯·æ±‚å‘é€åˆ°è‡ªèº«ï¼Œå¯èƒ½å¯¼è‡´ 404ã€‚ä¸ºé¿å…æ­¤ç±»é—®é¢˜ï¼š
	1. åœ¨ `frontend/.env.local` ä¸­è®¾ç½® `NEXT_PUBLIC_API_URL=http://localhost:8080`ï¼ˆæˆ–ä½ çš„åç«¯åœ°å€ï¼‰ï¼Œç„¶åé‡å¯ dev serverï¼›
	2. é¡¹ç›®å·²åœ¨ `next.config.ts` ä¸­ä¸º dev æ¨¡å¼æ·»åŠ äº† `rewrites`ï¼Œä¼šå°† `/api/v1/:path*` ä»£ç†åˆ° `NEXT_PUBLIC_API_URL`ï¼ˆå¦‚æœæœªè®¾ç½®åˆ™é»˜è®¤ `http://localhost:8080`ï¼‰ã€‚

ç¤ºä¾‹ `.env.local`ï¼š
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

é‡å¯ Dev Serverï¼š
```bash
cd frontend
npm run dev
```

## æ„å»ºä¸å¯¼å‡º

```bash
npm run build   # next.config.ts å·²é…ç½® output: 'export'ï¼Œäº§ç‰©è¾“å‡ºåˆ° frontend/out
```

æ„å»ºå®Œæˆåï¼Œ`frontend/out/` å¯ç›´æ¥ç”±åç«¯é™æ€æ‰˜ç®¡ï¼ˆå‚è§ `backend/internal/app/http_server/server.go`ï¼‰ã€‚

## æµ‹è¯•

```bash
npm test -- --coverage
```

æ ¸å¿ƒè¦†ç›–ï¼š

- `lib/api` æ‹¦æˆªå™¨ä¸é‰´æƒæµç¨‹
- `contexts/AuthContext` ç™»å½•/æ³¨é”€é€»è¾‘
- ç»„ä»¶ä¸é¡µé¢é›†æˆæµ‹è¯•ï¼ˆè§ `tests/`ï¼‰

## ä»£ç è§„èŒƒ

- ESLint + TypeScript ä¸¥æ ¼æ¨¡å¼
- SWR å…¨å±€ç¼“å­˜ï¼š`revalidateOnFocus=false`ï¼Œ`dedupingInterval=60s`ï¼ˆè§ `app/(protected)/layout.tsx`ï¼‰
- Prism.js è¯­è¨€åŒ…æŒ‰éœ€åŠ¨æ€åŠ è½½ï¼Œå‡å°‘é¦–åŒ…ä½“ç§¯

## ²âÑéÏà¹Ø×é¼ş

### ºËĞÄ×é¼ş

#### QuizQuestionCard£¨ÌâÄ¿¿¨Æ¬×é¼ş£©
- **Î»ÖÃ**£ºcomponents/quiz/QuizQuestionCard.tsx
- **¹¦ÄÜ**£ºÕ¹Ê¾µ¥¸ö²âÑéÌâÄ¿£¬Ö§³Öµ¥Ñ¡¡¢¶àÑ¡¡¢ÅĞ¶ÏÌâµÈ¶àÖÖÌâĞÍ
- **ÌØĞÔ**£º
  - ÎÈ¶¨µÄ A-D ±êÇ©äÖÈ¾£¨»ùÓÚÑ¡ÏîÊı×éË÷Òı£©
  - Ö§³Ö¶àÑ¡Ìâ/µ¥Ñ¡ÌâÇĞ»»
  - Ö§³Ö´úÂëÆ¬¶ÎÕ¹Ê¾
  - »Ø¹ËÄ£Ê½ÏÂÖ§³Ö½ûÓÃÑ¡Ôñ£¨disabled ÊôĞÔ£©
  - ÌâĞÍºÍÄÑ¶È±êÇ©ÏÔÊ¾

#### SubmitConfirmModal£¨Ìá½»È·ÈÏµ¯´°£©
- **Î»ÖÃ**£ºcomponents/quiz/SubmitConfirmModal.tsx
- **¹¦ÄÜ**£ºÔÚÓÃ»§Ìá½»Ç°ÏÔÊ¾ÒÑ´ğ/Î´´ğÌâÄ¿Í³¼Æ
- **ÌØĞÔ**£º
  - ÏÔÊ¾ÒÑ»Ø´ğºÍÎ´»Ø´ğÌâÄ¿ÊıÁ¿
  - ·ÀÖ¹ÓÃ»§Îó´¥Ìá½»°´Å¥
  - ÌáÊ¾ÓÃ»§¼ì²éÎ´´ğÌâÄ¿

#### QuizResultPage£¨½á¹ûÕ¹Ê¾Ò³£©
- **Î»ÖÃ**£ºcomponents/quiz/QuizResultPage.tsx
- **¹¦ÄÜ**£ºÕ¹Ê¾²âÑé½á¹û£¬°üÀ¨°Ù·ÖÖÆµÃ·ÖºÍÍ¨¹ı/Î´Í¨¹ı×´Ì¬
- **ÌØĞÔ**£º
  - °Ù·ÖÖÆµÃ·ÖÏÔÊ¾£¨0-100£©
  - Í¨¹ı×´Ì¬ÑÕÉ«Çø·Ö£¨ÂÌÉ«Í¨¹ı/ºìÉ«Î´Í¨¹ı£©
  - ÕıÈ·/´íÎóÌâÊıÍ³¼Æ
  - ²Ù×÷°´Å¥£ºÖØĞÂ²âÑé¡¢²é¿´½âÎö

#### QuestionTypeTag£¨ÌâĞÍ±êÇ©£©
- **Î»ÖÃ**£ºcomponents/quiz/QuestionTypeTag.tsx
- **¹¦ÄÜ**£ºÏÔÊ¾ÌâÄ¿ÀàĞÍ±êÇ©
- **ÌâĞÍÖ§³Ö**£ºµ¥Ñ¡Ìâ¡¢¶àÑ¡Ìâ¡¢¸Ä´íÌâ¡¢ÅĞ¶ÏÌâ¡¢´úÂëÊä³öÌâµÈ

#### QuizSkeletonLoader£¨¹Ç¼ÜÆÁ¼ÓÔØÆ÷£©
- **Î»ÖÃ**£ºcomponents/quiz/QuizSkeletonLoader.tsx
- **¹¦ÄÜ**£ºÌâÄ¿¼ÓÔØÊ±µÄÕ¼Î»·û¶¯»­
- **ÌØĞÔ**£ºÄ£ÄâÌâ¸É¡¢Ñ¡Ïî¡¢°´Å¥µÄ¼ÓÔØ¹ı³Ì

#### QuizMetaInfo£¨ÔªÊı¾İÕ¹Ê¾£©
- **Î»ÖÃ**£ºcomponents/quiz/QuizMetaInfo.tsx
- **¹¦ÄÜ**£ºÏÔÊ¾Ìâ¿âµÄÔªÊı¾İĞÅÏ¢
- **ÌØĞÔ**£º
  - ×ÜÌâÁ¿ÏÔÊ¾
  - Ô¤¼ÆÓÃÊ±¼ÆËã
  - ÄÑ¶È·Ö²¼Õ¹Ê¾
  - ÌâĞÍÕ¼±ÈÍ³¼Æ

#### AnswerIndicator£¨´ğ°¸Ö¸Ê¾Æ÷£©
- **Î»ÖÃ**£ºcomponents/quiz/AnswerIndicator.tsx
- **¹¦ÄÜ**£ºÔÚ»Ø¹ËÄ£Ê½ÖĞ¶Ô±ÈÓÃ»§´ğ°¸ÓëÕıÈ·´ğ°¸
- **ÌØĞÔ**£º
  - ÏÔÊ¾ÓÃ»§Ñ¡Ôñ
  - ÏÔÊ¾ÕıÈ·´ğ°¸
  - ´ğ¶Ô/´ğ´í×´Ì¬Ö¸Ê¾
  - ½âÎöĞÅÏ¢Õ¹Ê¾

### Ò³Ãæ×é¼ş

#### QuizPageClient£¨²âÑéÖ÷Ò³Ãæ£©
- **Î»ÖÃ**£ºpp/(protected)/quiz/[topic]/QuizPageClient.tsx
- **¹¦ÄÜ**£ºÍêÕûµÄ²âÑé½»»¥Á÷³Ì
- **ÌØĞÔ**£º
  - ÕÂ½ÚÑ¡ÔñºÍÌâÄ¿¼ÓÔØ
  - ´ğÌâµ¼º½ºÍ½ø¶ÈÕ¹Ê¾
  - Ìá½»È·ÈÏºÍ½á¹ûÕ¹Ê¾
  - ¹Ç¼ÜÆÁ¼ÓÔØÌ¬

#### QuizHistoryPage£¨ÀúÊ·ÁĞ±íÒ³£©
- **Î»ÖÃ**£ºpp/(protected)/quiz/history/page.tsx
- **¹¦ÄÜ**£ºÕ¹Ê¾ÓÃ»§µÄ²âÑéÀúÊ·¼ÇÂ¼
- **ÌØĞÔ**£º
  - °´Ê±¼äµ¹ĞòÅÅÁĞ
  - ÏÔÊ¾µÃ·Ö¡¢Í¨¹ı×´Ì¬¡¢Íê³ÉÊ±¼ä
  - Ö§³Öµã»÷½øÈëÏêÇé»Ø¹Ë

#### QuizReviewPage£¨»Ø¹ËÏêÇéÒ³£©
- **Î»ÖÃ**£ºpp/(protected)/quiz/history/[sessionId]/page.tsx
- **¹¦ÄÜ**£ºÕ¹Ê¾µ¥´Î²âÑéµÄÏêÏ¸»Ø¹Ë
- **ÌØĞÔ**£º
  - ÖğÌâÕ¹Ê¾´ğ°¸¶Ô±È
  - ÏÔÊ¾ÌâÄ¿½âÎö
  - Ö§³Ö´íÌâÉ¸Ñ¡ºÍÍ³¼Æ
  - ½ûÓÃÑ¡Ôñ½»»¥£¨Ö»¶ÁÄ£Ê½£©

### ×Ô¶¨Òå Hooks

#### useQuiz£¨²âÑéÂß¼­ Hook£©
- **Î»ÖÃ**£ºhooks/useQuiz.ts
- **¹¦ÄÜ**£º¹ÜÀíµ¥´Î²âÑéµÄÍêÕûÉúÃüÖÜÆÚ
- **¹ÜÀíµÄ×´Ì¬**£º
  - µ±Ç°ÌâÄ¿Ë÷Òı
  - ÓÃ»§´ğ°¸¼ÇÂ¼
  - Ìá½»½á¹û
  - ¼ÓÔØºÍ´íÎó×´Ì¬

#### useQuizHistory£¨ÀúÊ·Êı¾İ Hook£©
- **Î»ÖÃ**£ºhooks/useQuiz.ts
- **¹¦ÄÜ**£ºÊ¹ÓÃ SWR »ñÈ¡ÓÃ»§µÄ²âÑéÀúÊ·
- **ÌØĞÔ**£º
  - ×Ô¶¯»º´æºÍÖØĞÂÑéÖ¤
  - Ö§³Ö°´Ö÷Ìâ¹ıÂË

### ·şÎñºÍ¹¤¾ß

#### quizService£¨²âÑé·şÎñ£©
- **Î»ÖÃ**£ºservices/quizService.ts
- **Ö÷Òªº¯Êı**£º
  - etchQuizSession(topic, chapter)£º»ñÈ¡²âÑéÌâÄ¿
  - submitQuiz(payload)£ºÌá½»´ğ°¸
  - useQuizHistory(topic?)£º»ñÈ¡ÀúÊ·¼ÇÂ¼
  - useQuizReview(sessionId)£º»ñÈ¡»Ø¹ËÏêÇé
  - etchQuizStats(topic, chapter)£º»ñÈ¡Ìâ¿âÍ³¼Æ

#### ÀàĞÍ¶¨Òå
- **Î»ÖÃ**£º	ypes/quiz.ts
- **ºËĞÄÀàĞÍ**£º
  - QuizQuestion£ºµ¥¸öÌâÄ¿
  - QuizSubmitRequest£ºÌá½»ÇëÇó
  - QuizSubmitResult£ºÌá½»½á¹û
  - QuizHistoryItem£ºÀúÊ·Ïî
  - QuizReviewDetail£º»Ø¹ËÏêÇé
