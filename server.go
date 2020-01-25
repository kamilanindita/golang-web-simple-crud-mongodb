package main

import (
    "context"
    "log"
    "net/http"
    "text/template"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    
    
	
)

var ctx = func() context.Context {
	return context.Background()
}()

type Buku struct {
    Id primitive.ObjectID `bson:"_id, omitempty"`
    Penulis  string `bson:"penulis"`
    Judul string    `bson:"judul"`
    Kota string    `bson:"kota"`
    Penerbit string `bson:"penerbit"`
    Tahun interface{} `bson:"tahun"`
}

func connect() (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("mywebsite_crud"), nil
}


func HandlerIndex(w http.ResponseWriter, r *http.Request) {
    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Index.html",
        "views/Footer.html",
    ))
    data:=""
    var error = tmp.ExecuteTemplate(w,"Index",data)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}


func HandlerBuku(w http.ResponseWriter, r *http.Request) {
    db, err := connect()
    if err != nil {
        log.Fatal(err.Error())
    }

    csr, err := db.Collection("buku").Find(ctx, bson.M{})
    if err != nil {
        log.Fatal(err.Error())
    }
    defer csr.Close(ctx)

    result := make([]Buku, 0)
    for csr.Next(ctx) {
        var row Buku
        err := csr.Decode(&row)
        if err != nil {
            log.Fatal(err.Error())
        }

        result = append(result, row)
    }

    if len(result) > 0 {
        var tmp = template.Must(template.ParseFiles(
            "views/Header.html",
            "views/Menu.html",
            "views/Buku.html",
            "views/Footer.html",
        ))
       
        var error = tmp.ExecuteTemplate(w,"Buku",result)
        if error != nil {
            http.Error(w, error.Error(), http.StatusInternalServerError)
        }
    }
}


func HandlerTamabah(w http.ResponseWriter, r *http.Request) {
    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Tambah.html",
        "views/Footer.html",
    ))
    data:=""
    var error = tmp.ExecuteTemplate(w,"Tambah",data)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}


func HandlerSave(w http.ResponseWriter, r *http.Request) {
    db, err := connect()
    if err != nil {
        log.Fatal(err.Error())
    }

    penulis := r.FormValue("penulis")
    judul := r.FormValue("judul")
    kota := r.FormValue("kota")
    penerbit := r.FormValue("penerbit")
    tahun := r.FormValue("tahun")

    objID:= primitive.NewObjectID()

    var data = Buku{objID, penulis, judul, kota, penerbit, tahun}

    _, err = db.Collection("buku").InsertOne(ctx, data)
    if err != nil {
        log.Fatal(err.Error())
    }

    
    http.Redirect(w, r, "/buku", 301)
}



func HandlerEdit(w http.ResponseWriter, r *http.Request) {
    db, err := connect()
    if err != nil {
        log.Fatal(err.Error())
    }
    docID := r.URL.Query().Get("id") //mengabil id parameter berbentuk hexsa
    objID, err := primitive.ObjectIDFromHex(docID)
    if err != nil {
        log.Fatal(err.Error())
    }
    csr, err := db.Collection("buku").Find(ctx, bson.M{ "_id": bson.M{"$eq": objID} })
    if err != nil {
        log.Fatal(err.Error())
    }
    defer csr.Close(ctx)

    result := make([]Buku, 0)
    for csr.Next(ctx) {
        var row Buku
        err := csr.Decode(&row)
        if err != nil {
            log.Fatal(err.Error())
        }

        result = append(result, row)
    }

    if len(result) > 0 {
    
    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Edit.html",
        "views/Footer.html",
    ))

    var error = tmp.ExecuteTemplate(w,"Edit",result[0])
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }

    }
}


func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
    db, err := connect()
    if err != nil {
        log.Fatal(err.Error())
    }
    penulis := r.FormValue("penulis")
    judul := r.FormValue("judul")
    kota := r.FormValue("kota")
    penerbit := r.FormValue("penerbit")
    tahun := r.FormValue("tahun")

    docID := r.URL.Query().Get("id") //mengabil id parameter berbentuk hexsa
    objID, err := primitive.ObjectIDFromHex(docID)
 
    var changes = Buku{objID, penulis, judul, kota, penerbit, tahun}
 
    _, err = db.Collection("buku").UpdateOne(ctx, bson.M{ "_id": bson.M{"$eq": objID} } , bson.M{"$set": changes} )
    if err != nil {
        log.Fatal(err.Error())
    }
    
    http.Redirect(w, r, "/buku", 301)
}


func HandlerDelete(w http.ResponseWriter, r *http.Request) {
    db, err := connect()
    if err != nil {
        log.Fatal(err.Error())
    }

    docID := r.URL.Query().Get("id") //mengabil id parameter berbentuk hexsa
    objID, err := primitive.ObjectIDFromHex(docID)
    if err != nil {
        log.Fatal(err.Error())
    }

    _, err = db.Collection("buku").DeleteOne(ctx, bson.M{ "_id": bson.M{"$eq": objID} })
    if err != nil {
        log.Fatal(err.Error())
    }
 
    http.Redirect(w, r, "/buku", 301)
}


func main() {
    log.Println("Server started on: http://localhost:8000")
    http.HandleFunc("/", HandlerIndex)
    http.HandleFunc("/buku", HandlerBuku)
    http.HandleFunc("/buku/tambah", HandlerTamabah)
    http.HandleFunc("/buku/edit", HandlerEdit)
    http.HandleFunc("/buku/save", HandlerSave)
    http.HandleFunc("/buku/update", HandlerUpdate)
    http.HandleFunc("/buku/delete", HandlerDelete)
    http.ListenAndServe(":8000", nil)
}