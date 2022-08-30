int val;
int tempPin = 1;

float input_volt = 0.0;
float temp=0.0;
float r1=1000.0;    //r1 value
float r2=10000.0;      //r2 value

void setup()
{
  Serial.begin(9600);
}

void loop()
{
    int analogvalue = analogRead(A2);
    temp = (analogvalue * 5.0) / 1024.0;  // FORMULA USED TO CONVERT THE VOLTAGE
    input_volt = temp / (r2/(r1+r2));
    if (input_volt < 0.1) 
    {
      input_volt=0.0;
    }                  // prints the voltage value in the serial monitor
    Serial.print(input_volt, 6);
    Serial.print(" ");

  
  val = analogRead(tempPin);
  float mv = ( val/1024.0)*5000;
  float cel = mv/10;
  float farh = (cel*9)/5 + 32;
  Serial.print(cel, 6);
  Serial.println();
  delay(100);
}