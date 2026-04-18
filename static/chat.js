
function scrollToBottom() {
	const chatDisplay = document.querySelector('.chat-display');
	console.log('chatDisplay:', chatDisplay);
	if (chatDisplay) {
		console.log('scrollTop:', chatDisplay.scrollTop, 'scrollHeight:', chatDisplay.scrollHeight);
		chatDisplay.scrollTop = chatDisplay.scrollHeight;
		console.log('new scrollTop:', chatDisplay.scrollTop);
	}
}

document.addEventListener('DOMContentLoaded', scrollToBottom);
if (document.readyState === 'complete' || document.readyState !== 'loading') {
	setTimeout(scrollToBottom, 0);
}

function saveInputState() {
	const textBox = document.getElementById('text');
	const nickBox = document.getElementById('nickname');
	sessionStorage.setItem('chat_text', textBox ? textBox.value : '');
	sessionStorage.setItem('chat_nickname', nickBox ? nickBox.value : '');
	sessionStorage.setItem('chat_focused', document.activeElement === textBox ? 'true' : 'false');
}

function restoreInputState() {
	const textBox = document.getElementById('text');
	const nickBox = document.getElementById('nickname');
	if (textBox) textBox.value = sessionStorage.getItem('chat_text') || '';
	if (nickBox) nickBox.value = sessionStorage.getItem('chat_nickname') || '';
	if (sessionStorage.getItem('chat_focused') === 'true' && textBox) {
		textBox.focus();
	}
}

async function fetchChatData() {
	let intervalId;
	intervalId = setInterval(async () => {
		try {
			saveInputState();
			location.reload();
		} catch (error) {
			console.error("Error during polling:", error);
		}
	}, 5000);
	scrollToBottom();
	restoreInputState();
}

document.addEventListener('DOMContentLoaded', function() {
  const textBox = document.getElementById('text');
  const sendBtn = document.getElementById('send');
  const nickBox = document.getElementById('nickname');

  if (!textBox || !sendBtn) return;

  async function sendMessage() {
    const text = textBox.value.trim();
    const nick = nickBox ? nickBox.value.trim() : 'anon';
    if (!text) return;

    const params = new URLSearchParams({ nickname: nick, text: text });
    await fetch('/send?' + params.toString());
    textBox.value = '';
    location.reload();
  }

  sendBtn.addEventListener('click', function(e) {
    e.preventDefault();
    sendMessage();
  });

  textBox.addEventListener('keydown', function(e) {
    if (e.key === 'Enter' && e.shiftKey) {
      e.preventDefault();
      const start = this.selectionStart;
      const end = this.selectionEnd;
      this.value = this.value.substring(0, start) + '\n' + this.value.substring(end);
      this.selectionStart = this.selectionEnd = start + 1;
    } else if (e.key === 'Enter') {
      e.preventDefault();
      sendMessage();
    }
  });
});

fetchChatData();




