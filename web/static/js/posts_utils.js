export function decodeHTMLEntities(text) {
    const textarea = document.createElement('textarea');
    textarea.innerHTML = text;
    return textarea.value;
}

export function formatTextWithLineBreaks(text, contentElement) {
   // Decode HTML entities
           const decodedContent = decodeHTMLEntities(text);
   
           // Split by newlines and add <br> elements
           const lines = decodedContent.split('\n');
           lines.forEach((line, index) => {
               contentElement.appendChild(document.createTextNode(line));
               if (index < lines.length - 1) {
                   contentElement.appendChild(document.createElement('br'));
               }
           });
}